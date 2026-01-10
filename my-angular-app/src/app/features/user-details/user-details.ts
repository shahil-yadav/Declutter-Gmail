import { Component, OnInit, isDevMode } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { User } from '../../core/models/user';
import { UsersService } from '../../core/services/users';

@Component({
  selector: 'app-user-details',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    RouterModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatSlideToggleModule
  ],
  templateUrl: './user-details.html',
  styleUrl: './user-details.scss',
})
export class UserDetailsComponent implements OnInit {
  user: User | null = null;
  loading = true;
  errorMessage: string | null = null;

  // Dev mode helpers
  isDev = isDevMode();
  simulateError = false;

  constructor(
    private route: ActivatedRoute,
    private usersService: UsersService
  ) { }

  ngOnInit(): void {
    this.loadCurrentUser();
  }

  loadCurrentUser(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.loadUser(id);
    } else {
      this.loading = false;
      this.errorMessage = 'Invalid User ID';
    }
  }

  loadUser(id: string): void {
    if (this.simulateError) {
      this.errorMessage = 'Simulated Network Error: Unable to connect to the backend.';
      this.user = null;
      this.loading = false;
      return;
    }

    this.loading = true;
    this.usersService.getUser(id).subscribe({
      next: (response) => {
        if (response.code === 200) {
          this.user = response.data;
          this.errorMessage = null;
        } else if (response.code === 404) {
          this.errorMessage = 'User not found! Please do not manually change the URL.';
        } else {
          this.errorMessage = response.msg || 'Failed to load user details';
        }
        this.loading = false;
      },
      error: (err) => {
        this.errorMessage = 'Unable to connect to the backend.';
        this.loading = false;
        console.error('Error fetching user', err);
      }
    });
  }

  toggleErrorSimulation(): void {
    this.loadCurrentUser();
  }
}
