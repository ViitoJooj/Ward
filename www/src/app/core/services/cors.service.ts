import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { CorsItem, CorsOutput } from '../../features/dashboard/dashboard.models';

@Injectable({ providedIn: 'root' })
export class CorsService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  getAll(): Observable<CorsItem[]> {
    return this.http.get<CorsItem[]>(`${this.apiUrl}/cors/`);
  }

  create(origin: string): Observable<CorsOutput> {
    return this.http.post<CorsOutput>(`${this.apiUrl}/cors/`, { origin });
  }

  update(id: number, origin: string): Observable<CorsOutput> {
    return this.http.put<CorsOutput>(`${this.apiUrl}/cors/${id}`, { origin });
  }

  delete(id: number): Observable<CorsOutput> {
    return this.http.delete<CorsOutput>(`${this.apiUrl}/cors/${id}`);
  }
}
