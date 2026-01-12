import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ApiResponse } from '../models/user';

@Injectable({
    providedIn: 'root'
})
export class AuthService {
    private apiUrl = 'http://localhost:8080/auth';

    constructor(private http: HttpClient) { }

    getLoginUrl(): Observable<ApiResponse<string>> {
        return this.http.get<ApiResponse<string>>(`${this.apiUrl}/login?redirect=false`, { withCredentials: true });
    }
}
