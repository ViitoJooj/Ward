import { Component, OnInit, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { RouteRuleService, RouteRuleInput } from '../../../../core/services/route-rule.service';
import { RouteRule } from '../../dashboard.models';

@Component({
  selector: 'app-routes',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './routes.component.html',
  styleUrl: './routes.component.scss'
})
export class RoutesComponent implements OnInit {
  rules = signal<RouteRule[]>([]);
  loading = signal(true);
  error = signal<string | null>(null);
  success = signal<string | null>(null);

  showCreate = signal(false);
  creating = signal(false);
  createError = signal<string | null>(null);

  editingId = signal<number | null>(null);
  savingId = signal<number | null>(null);
  deletingId = signal<number | null>(null);

  pathFilter = signal('');

  filteredRules = computed(() => {
    const f = this.pathFilter().toLowerCase().trim();
    if (!f) return this.rules();
    return this.rules().filter(r =>
      r.path.toLowerCase().includes(f) || r.method.toLowerCase().includes(f)
    );
  });

  createForm: FormGroup;
  editForm: FormGroup;

  constructor(private service: RouteRuleService, private fb: FormBuilder) {
    this.createForm = this.buildForm();
    this.editForm = this.buildForm();
  }

  private buildForm(): FormGroup {
    return this.fb.group({
      path: ['', Validators.required],
      method: [''],
      rate_limit_enabled: [false],
      rate_limit_rps: [10, [Validators.min(0.0001)]],
      rate_limit_burst: [20, [Validators.min(1)]],
      target_url: [''],
      geo_routing_enabled: [false],
      enabled: [true],
    });
  }

  ngOnInit(): void {
    this.load();
  }

  load(): void {
    this.loading.set(true);
    this.error.set(null);
    this.service.getAll().subscribe({
      next: (data) => {
        this.rules.set(data);
        this.loading.set(false);
      },
      error: () => {
        this.error.set('Failed to load route rules.');
        this.loading.set(false);
      }
    });
  }

  toggleCreate(): void {
    this.showCreate.update(v => !v);
    this.createError.set(null);
    this.createForm.reset({
      path: '', method: '', rate_limit_enabled: false,
      rate_limit_rps: 10, rate_limit_burst: 20,
      target_url: '', geo_routing_enabled: false, enabled: true,
    });
  }

  onCreate(): void {
    if (this.createForm.invalid || this.creating()) return;
    this.creating.set(true);
    this.createError.set(null);
    this.service.create(this.toInput(this.createForm)).subscribe({
      next: () => {
        this.creating.set(false);
        this.showCreate.set(false);
        this.success.set('Route rule created.');
        this.load();
        this.clearSuccess();
      },
      error: (err) => {
        this.creating.set(false);
        this.createError.set(err?.error?.message ?? 'Failed to create route rule.');
      }
    });
  }

  startEdit(rule: RouteRule): void {
    this.editingId.set(rule.id);
    this.editForm.patchValue({
      path: rule.path,
      method: rule.method,
      rate_limit_enabled: rule.rate_limit_enabled,
      rate_limit_rps: rule.rate_limit_rps,
      rate_limit_burst: rule.rate_limit_burst,
      target_url: rule.target_url,
      geo_routing_enabled: rule.geo_routing_enabled,
      enabled: rule.enabled,
    });
  }

  cancelEdit(): void {
    this.editingId.set(null);
  }

  saveEdit(rule: RouteRule): void {
    if (this.savingId() !== null) return;
    this.savingId.set(rule.id);
    this.error.set(null);
    this.service.update(rule.id, this.toInput(this.editForm)).subscribe({
      next: () => {
        this.savingId.set(null);
        this.editingId.set(null);
        this.success.set('Route rule updated.');
        this.load();
        this.clearSuccess();
      },
      error: (err) => {
        this.savingId.set(null);
        this.error.set(err?.error?.message ?? 'Failed to update route rule.');
      }
    });
  }

  deleteRule(rule: RouteRule): void {
    this.deletingId.set(rule.id);
    this.error.set(null);
    this.service.delete(rule.id).subscribe({
      next: () => {
        this.deletingId.set(null);
        this.success.set('Route rule deleted.');
        this.load();
        this.clearSuccess();
      },
      error: (err) => {
        this.deletingId.set(null);
        this.error.set(err?.error?.message ?? 'Failed to delete route rule.');
      }
    });
  }

  onFilterInput(event: Event): void {
    this.pathFilter.set((event.target as HTMLInputElement).value);
  }

  isRateLimitEnabled(): boolean {
    return !!this.editForm.get('rate_limit_enabled')?.value;
  }

  isCreateRateLimitEnabled(): boolean {
    return !!this.createForm.get('rate_limit_enabled')?.value;
  }

  private toInput(form: FormGroup): RouteRuleInput {
    const v = form.value;
    return {
      path: v.path,
      method: v.method ?? '',
      rate_limit_enabled: !!v.rate_limit_enabled,
      rate_limit_rps: Number(v.rate_limit_rps),
      rate_limit_burst: Number(v.rate_limit_burst),
      target_url: v.target_url ?? '',
      geo_routing_enabled: !!v.geo_routing_enabled,
      enabled: !!v.enabled,
    };
  }

  private clearSuccess(): void {
    setTimeout(() => this.success.set(null), 3000);
  }
}
