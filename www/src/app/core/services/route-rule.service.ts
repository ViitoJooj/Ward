import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from '../../../environments/environment';
import { RouteRule, RouteRuleListOutput, RouteRuleOutput } from '../../features/dashboard/dashboard.models';

export interface RouteRuleInput {
  path: string;
  method: string;
  rate_limit_enabled: boolean;
  rate_limit_rps: number;
  rate_limit_burst: number;
  target_url: string;
  geo_routing_enabled: boolean;
  enabled: boolean;
}

@Injectable({ providedIn: 'root' })
export class RouteRuleService {
  private readonly base = `${environment.apiUrl}/route-rules`;

  constructor(private http: HttpClient) {}

  getAll(): Observable<RouteRule[]> {
    return this.http.get<RouteRuleListOutput>(this.base).pipe(map(r => r.data ?? []));
  }

  create(input: RouteRuleInput): Observable<RouteRuleOutput> {
    return this.http.post<RouteRuleOutput>(this.base, input);
  }

  update(id: number, input: RouteRuleInput): Observable<RouteRuleOutput> {
    return this.http.put<RouteRuleOutput>(`${this.base}/${id}`, input);
  }

  delete(id: number): Observable<RouteRuleOutput> {
    return this.http.delete<RouteRuleOutput>(`${this.base}/${id}`);
  }
}
