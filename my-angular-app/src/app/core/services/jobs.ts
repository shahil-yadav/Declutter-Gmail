import { Injectable } from '@angular/core';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ApiResponse } from '../models/user';

@Injectable({
    providedIn: 'root'
})
export class JobsService {
    private apiUrl = 'http://localhost:8080/v1/job';

    constructor(private http: HttpClient) { }

    startScan(userId: string): Observable<ApiResponse<string>> {
        const formData = new FormData();
        formData.append('user-id', userId);
        return this.http.post<ApiResponse<string>>(`${this.apiUrl}/scan`, formData);
    }

    startTrashJob(userId: string, senders: string[]): Observable<ApiResponse<string>> {
        let params = new HttpParams();
        params = params.append('user-id', userId);
        senders.forEach(sender => {
            params = params.append('sender[]', sender);
        });

        const headers = new HttpHeaders({ 'Content-Type': 'application/x-www-form-urlencoded' });

        return this.http.post<ApiResponse<string>>(`${this.apiUrl}/trash`, params.toString(), { headers });
    }
}
