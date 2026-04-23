import { Component, OnInit, signal, computed, HostListener, ElementRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { ApplicationService } from '../../../../core/services/application.service';
import { Application } from '../../dashboard.models';

export const COUNTRIES = [
  { code: 'AR', name: 'Argentina' },
  { code: 'AU', name: 'Australia' },
  { code: 'BR', name: 'Brazil' },
  { code: 'CA', name: 'Canada' },
  { code: 'CN', name: 'China' },
  { code: 'FR', name: 'France' },
  { code: 'DE', name: 'Germany' },
  { code: 'IN', name: 'India' },
  { code: 'IT', name: 'Italy' },
  { code: 'JP', name: 'Japan' },
  { code: 'MX', name: 'Mexico' },
  { code: 'PT', name: 'Portugal' },
  { code: 'RU', name: 'Russia' },
  { code: 'ZA', name: 'South Africa' },
  { code: 'ES', name: 'Spain' },
  { code: 'GB', name: 'United Kingdom' },
  { code: 'US', name: 'United States' },
];

@Component({
  selector: 'app-applications',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './applications.component.html',
  styleUrl: './applications.component.scss'
})
export class ApplicationsComponent implements OnInit {
  applications = signal<Application[]>([]);
  loading = signal(true);
  errorMessage = signal<string | null>(null);
  showCreateForm = signal(false);
  creating = signal(false);
  createError = signal<string | null>(null);
  deletingId = signal<number | null>(null);

  countrySearchTerm = signal('');
  isCountryDropdownOpen = signal(false);

  filteredCountries = computed(() => {
    const term = this.countrySearchTerm().toLowerCase().trim();
    if (!term) return COUNTRIES;
    return COUNTRIES.filter(c => c.name.toLowerCase().includes(term) || c.code.toLowerCase().includes(term));
  });

  createForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    private applicationService: ApplicationService,
    private eRef: ElementRef
  ) {
    this.createForm = this.fb.group({
      url: ['', Validators.required],
      country: ['', Validators.required]
    });
  }

  ngOnInit(): void {
    this.loadApplications();
  }

  loadApplications(): void {
    this.loading.set(true);
    this.errorMessage.set(null);

    this.applicationService.getAll().subscribe({
      next: (data) => {
        this.applications.set(data ?? []);
        this.loading.set(false);
      },
      error: () => {
        this.errorMessage.set('Erro ao carregar aplicações.');
        this.loading.set(false);
      }
    });
  }

  toggleCreateForm(): void {
    this.showCreateForm.update(v => !v);
    this.createError.set(null);
    this.createForm.reset({ url: '', country: '' });
    this.countrySearchTerm.set('');
    this.isCountryDropdownOpen.set(false);
  }

  @HostListener('document:click', ['$event'])
  clickout(event: Event) {
    if (!this.eRef.nativeElement.contains(event.target)) {
      this.isCountryDropdownOpen.set(false);
    }
  }

  onCountrySearch(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.countrySearchTerm.set(input.value);
    this.isCountryDropdownOpen.set(true);
  }

  selectCountry(code: string, name: string): void {
    this.createForm.patchValue({ country: code });
    this.countrySearchTerm.set(`${name} (${code})`);
    this.isCountryDropdownOpen.set(false);
  }

  openCountryDropdown(): void {
    this.isCountryDropdownOpen.set(true);
    if (!this.createForm.get('country')?.value) {
      this.countrySearchTerm.set('');
    }
  }

  onCreate(): void {
    if (this.createForm.invalid) return;

    this.creating.set(true);
    this.createError.set(null);

    this.applicationService.create(this.createForm.value).subscribe({
      next: () => {
        this.creating.set(false);
        this.showCreateForm.set(false);
        this.createForm.reset();
        this.loadApplications();
      },
      error: () => {
        this.creating.set(false);
        this.createError.set('Erro ao criar aplicação. Verifique os dados e tente novamente.');
      }
    });
  }

  onDelete(id: number): void {
    this.deletingId.set(id);

    this.applicationService.delete(id).subscribe({
      next: () => {
        this.deletingId.set(null);
        this.loadApplications();
      },
      error: () => {
        this.deletingId.set(null);
        this.errorMessage.set('Erro ao remover aplicação.');
      }
    });
  }
}
