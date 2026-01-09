import { Component, inject } from '@angular/core';
import { UsersService } from '../users.service';
import { User } from '../types';
import { RouterLink } from '@angular/router';
@Component({
  selector: 'app-home',
  templateUrl: './home.page.html',
  styleUrl: './home.page.css',
  imports: [RouterLink],
})
export class HomePage {
  private userService = inject(UsersService);
  networkError = '';
  users: User[] = [];

  constructor() {
    this.userService
      .getUsers()
      .then((users) => (this.users = users.data))
      .catch((err) => (this.networkError = err));
  }
}
