import { Component, OnDestroy, OnInit, computed, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { forkJoin, Subscription } from 'rxjs';
import { HealthService } from '../../../../core/services/health.service';
import { HealthOverviewData, HealthRouteData, HealthStatus } from '../../dashboard.models';

const ROUTE_LIMIT = 50;
const AUTO_REFRESH_SECONDS = 30;
const DONUT_R = 36;
const DONUT_CIRC = 2 * Math.PI * DONUT_R;

interface DonutSegment {
  label: string;
  count: number;
  pct: number;
  color: string;
  dasharray: string;
  dashoffset: number;
}

@Component({
  selector: 'app-health',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './health.component.html',
  styleUrl: './health.component.scss'
})
export class HealthComponent implements OnInit, OnDestroy {
  loading = signal(true);
  error = signal<string | null>(null);
  windowMinutes = signal(15);
  overview = signal<HealthOverviewData | null>(null);
  routes = signal<HealthRouteData[]>([]);
  countdown = signal(AUTO_REFRESH_SECONDS);

  routeFilter = signal('');
  statusFilter = signal<HealthStatus | 'all'>('all');
  responseFilter = signal<'all' | '2xx' | '4xx' | '5xx'>('all');

  filteredRoutes = computed(() => {
    const rf = this.routeFilter().toLowerCase().trim();
    const sf = this.statusFilter();
    const resp = this.responseFilter();
    return this.routes().filter(r => {
      if (rf && !r.path.toLowerCase().includes(rf) && !r.method.toLowerCase().includes(rf)) return false;
      if (sf !== 'all' && r.status !== sf) return false;
      if (resp === '5xx' && r.server_error_rate === 0) return false;
      if (resp === '4xx' && r.client_error_rate === 0) return false;
      if (resp === '2xx' && (r.server_error_rate > 0 || r.client_error_rate > 0)) return false;
      return true;
    });
  });

  topRoutes = computed(() =>
    [...this.routes()].sort((a, b) => b.request_count - a.request_count).slice(0, 6)
  );

  maxRequestCount = computed(() =>
    Math.max(...this.routes().map(r => r.request_count), 1)
  );

  donutSegments = computed((): DonutSegment[] => {
    const o = this.overview();
    if (!o || o.total_requests === 0) return [];
    const total = o.total_requests;
    const rawSegs = [
      { label: '2xx', count: Math.max(0, total - o.server_errors - o.client_errors), color: 'var(--ward-healthy)' },
      { label: '4xx', count: o.client_errors, color: 'var(--ward-degraded)' },
      { label: '5xx', count: o.server_errors, color: 'var(--ward-unhealthy)' },
    ].filter(s => s.count > 0);
    let acc = 0;
    return rawSegs.map(seg => {
      const len = (seg.count / total) * DONUT_CIRC;
      const item: DonutSegment = {
        ...seg,
        pct: (seg.count / total) * 100,
        dasharray: `${len.toFixed(2)} ${DONUT_CIRC.toFixed(2)}`,
        dashoffset: -acc,
      };
      acc += len;
      return item;
    });
  });

  private pollSubscription?: Subscription;
  private countdownTimer?: ReturnType<typeof setInterval>;

  constructor(private healthService: HealthService) {}

  ngOnInit(): void {
    this.loadHealth();
    this.startCountdown();
  }

  ngOnDestroy(): void {
    this.pollSubscription?.unsubscribe();
    clearInterval(this.countdownTimer);
  }

  private startCountdown(): void {
    clearInterval(this.countdownTimer);
    this.countdown.set(AUTO_REFRESH_SECONDS);
    this.countdownTimer = setInterval(() => {
      const next = this.countdown() - 1;
      if (next <= 0) {
        this.countdown.set(AUTO_REFRESH_SECONDS);
        this.loadHealth();
      } else {
        this.countdown.set(next);
      }
    }, 1000);
  }

  reload(): void {
    this.startCountdown();
    this.loadHealth();
  }

  setWindow(minutes: number): void {
    if (this.windowMinutes() === minutes) return;
    this.windowMinutes.set(minutes);
    this.reload();
  }

  setStatusFilter(f: HealthStatus | 'all'): void { this.statusFilter.set(f); }
  setResponseFilter(f: 'all' | '2xx' | '4xx' | '5xx'): void { this.responseFilter.set(f); }

  routeBarWidth(route: HealthRouteData): number {
    return (route.request_count / this.maxRequestCount()) * 100;
  }

  statusClass(status: HealthStatus): string {
    return `status-badge status-badge--${status}`;
  }

  metricPercent(value: number): string { return `${value.toFixed(2)}%`; }
  metricLatency(value: number): string { return `${value.toFixed(0)} ms`; }
  metricRpm(value: number): string { return value.toFixed(2); }

  availabilityBar(): number { return Math.min(100, Math.max(0, this.overview()?.availability ?? 0)); }
  serverErrorBar(): number { return Math.min(100, ((this.overview()?.server_error_rate ?? 0) / 10) * 100); }
  clientErrorBar(): number { return Math.min(100, ((this.overview()?.client_error_rate ?? 0) / 20) * 100); }
  avgLatencyBar(): number { return Math.min(100, ((this.overview()?.average_latency_ms ?? 0) / 2000) * 100); }
  p95LatencyBar(): number { return Math.min(100, ((this.overview()?.p95_latency_ms ?? 0) / 4000) * 100); }

  latencyColor(ms: number): string {
    if (ms >= 1500) return 'unhealthy';
    if (ms >= 700) return 'degraded';
    return 'healthy';
  }

  errorRateColor(rate: number): string {
    if (rate >= 10) return 'unhealthy';
    if (rate >= 3) return 'degraded';
    return 'healthy';
  }

  availabilityColor(avail: number): string {
    if (avail < 95) return 'unhealthy';
    if (avail < 98) return 'degraded';
    return 'healthy';
  }

  private loadHealth(): void {
    this.pollSubscription?.unsubscribe();
    this.loading.set(true);
    this.error.set(null);
    this.pollSubscription = forkJoin({
      overview: this.healthService.getOverview(this.windowMinutes()),
      routes: this.healthService.getRoutes(this.windowMinutes(), ROUTE_LIMIT),
    }).subscribe({
      next: (res) => {
        this.overview.set(res.overview.data);
        this.routes.set(res.routes.data ?? []);
        this.loading.set(false);
      },
      error: (err) => {
        this.error.set(err?.error?.message ?? 'Failed to load health data.');
        this.loading.set(false);
      }
    });
  }
}
