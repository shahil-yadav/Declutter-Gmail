import { Injectable } from '@angular/core';
import axios from 'axios';
import { Response, User } from './types';

type UsersResponse = Response<User[]>;
type UserResponse = Response<User>;

@Injectable({
  providedIn: 'root',
})
export class UsersService {
  private url = `http://127.0.0.1:7331/v1`;

  async getUsers(): Promise<UsersResponse> {
    const resp = await axios.get<UsersResponse>(`${this.url}/users`);
    return resp.data;
  }

  async getUser(userId: string): Promise<UserResponse> {
    const res = await axios.get<UserResponse>(`${this.url}/users/${userId}`);
    return res.data;
  }
}
