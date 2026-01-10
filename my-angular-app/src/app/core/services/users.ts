import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ApiResponse, User } from '../models/user';

@Injectable({
    providedIn: 'root'
})
export class UsersService {
    private apiUrl = 'http://127.0.0.1:7331/v1/users';

    constructor(private http: HttpClient) { }

    getUsers(): Observable<ApiResponse<User[]>> {
        return this.http.get<ApiResponse<User[]>>(this.apiUrl);
    }

    getUser(id: string): Observable<ApiResponse<User>> {
        return this.http.get<ApiResponse<User>>(`${this.apiUrl}/${id}`);
    }
}
