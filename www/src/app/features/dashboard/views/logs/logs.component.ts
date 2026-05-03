import { Component, OnInit, OnDestroy, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Subscription, interval, switchMap, startWith } from 'rxjs';
import { RequestLogService } from '../../../../core/services/request-log.service';
import { RequestLog } from '../../dashboard.models';

const POLL_INTERVAL_MS = 3000;

@Component({
  selector: 'app-logs',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './logs.component.html',
  styleUrl: './logs.component.scss'
})
export class LogsComponent implements OnInit, OnDestroy {
  logs = signal<RequestLog[]>([]);
  loading = signal(true);
  expandedId = signal<number | null>(null);
  filterPath = signal('');
  filterStatus = signal<'all' | '2xx' | '3xx' | '4xx' | '5xx'>('all');
  readonly statusOptions = ['all', '2xx', '3xx', '4xx', '5xx'] as const;

  filtered = computed(() => {
    const path = this.filterPath().toLowerCase();
    const status = this.filterStatus();
    return this.logs().filter(log => {
      const pathMatch = !path || log.path.toLowerCase().includes(path);
      const statusMatch = status === 'all'
        || (status === '2xx' && log.status_code >= 200 && log.status_code < 300)
        || (status === '3xx' && log.status_code >= 300 && log.status_code < 400)
        || (status === '4xx' && log.status_code >= 400 && log.status_code < 500)
        || (status === '5xx' && log.status_code >= 500);
      return pathMatch && statusMatch;
    });
  });

  private pollSubscription?: Subscription;

  constructor(private requestLogService: RequestLogService) {}

  ngOnInit(): void {
    this.pollSubscription = interval(POLL_INTERVAL_MS).pipe(
      startWith(0),
      switchMap(() => this.requestLogService.getAll())
    ).subscribe({
      next: (res) => {
        const visible = (res.data ?? []).filter(log => log.path !== '/ward/api/v1/logs');
        this.logs.set(visible);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
      }
    });
  }

  ngOnDestroy(): void {
    this.pollSubscription?.unsubscribe();
  }

  toggleRow(id: number): void {
    this.expandedId.set(this.expandedId() === id ? null : id);
  }

  onFilterPathInput(event: Event): void {
    this.filterPath.set((event.target as HTMLInputElement).value);
  }

  setFilterStatus(s: 'all' | '2xx' | '3xx' | '4xx' | '5xx'): void {
    this.filterStatus.set(s);
  }

  methodClass(method: string): string {
    return `method method--${method.toLowerCase()}`;
  }

  statusClass(code: number): string {
    if (code < 300) return 'status status--2xx';
    if (code < 400) return 'status status--3xx';
    if (code < 500) return 'status status--4xx';
    return 'status status--5xx';
  }

  formatBytes(bytes: number): string {
    if (bytes === 0) return '0 B';
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / 1024 / 1024).toFixed(1)} MB`;
  }
}
