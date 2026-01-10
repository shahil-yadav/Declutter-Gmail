import { Component, OnInit, isDevMode } from '@angular/core';
import { RouterModule } from '@angular/router';

import { FormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { User } from '../../core/models/user';
import { UsersService } from '../../core/services/users';
import { RelativeTimePipe } from '../../shared/pipes/relative-time-pipe';
import { AuthService } from '../../core/services/auth';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    FormsModule,
    RouterModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatSlideToggleModule,
    RelativeTimePipe
  ],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.scss',
})
export class DashboardComponent implements OnInit {
  users: User[] = [];
  errorMessage: string | null = null;
  loginUrl: string | null = null;

  // Dev mode helpers
  isDev = isDevMode();
  simulateError = false;

  constructor(
    private usersService: UsersService,
    private authService: AuthService
  ) { }

  ngOnInit(): void {
    this.loadUsers();
    this.fetchLoginUrl();
  }

  loadUsers(): void {
    if (this.simulateError) {
      this.errorMessage = 'Simulated Network Error: Unable to connect to the backend.';
      this.users = [];
      return;
    }

    this.usersService.getUsers().subscribe({
      next: (response) => {
        if (response.code === 200) {
          this.users = response.data || [];
          this.errorMessage = null;
        } else {
          this.errorMessage = response.msg || 'Failed to load users';
        }
      },
      error: (err) => {
        this.errorMessage = 'Unable to connect to the backend. Please try again later.';
        console.error('Error fetching users', err);
      },
    });
  }

  fetchLoginUrl(): void {
    this.authService.getLoginUrl().subscribe({
      next: (response) => {
        if (response.code === 200 && response.data) {
          this.loginUrl = response.data;
        } else {
          console.error('Failed to get login URL', response.msg);
        }
      },
      error: (err) => {
        console.error('Error fetching login URL', err);
      }
    });
  }

  toggleErrorSimulation(): void {
    this.loadUsers();
  }
}
