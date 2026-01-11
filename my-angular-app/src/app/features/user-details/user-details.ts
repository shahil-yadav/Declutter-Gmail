import { Component, OnInit, isDevMode } from '@angular/core';
import { Observable } from 'rxjs'; // Fixed import
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { User, ScanStatus, ScanResult } from '../../core/models/user';
import { UsersService } from '../../core/services/users';

import { JobsService } from '../../core/services/jobs';

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
  scanMessage: string | null = null;
  scanStatus: ScanStatus | null = null;

  // Dev mode helpers
  isDev = isDevMode();
  simulateError = false;
  simulatedScanSuccess = false;
  simulatedScanError = false;
  scanErrorMessage: string | null = null;

  // Selection & Modal
  selectedSenders = new Set<ScanResult>();
  showDeleteModal = false;

  get totalDeleteCount(): number {
    let count = 0;
    this.selectedSenders.forEach(s => count += s.Count);
    return count;
  }

  get selectedSendersList(): ScanResult[] {
    return Array.from(this.selectedSenders);
  }

  get isAllSelected(): boolean {
    const results = this.scanStatus?.Results || [];
    return results.length > 0 && this.selectedSenders.size === results.length;
  }

  constructor(
    private route: ActivatedRoute,
    private usersService: UsersService,
    private jobsService: JobsService
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
          this.loadScanStatus(this.user.UserId).subscribe(() => {
            // Only load trash status after scan status (or parallel, doesnt matter much but safer)
            // simplified: parallel is fine
          });
          this.loadScanStatus(this.user.UserId).subscribe();
          this.loadTrashStatus(this.user.UserId);
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

  loadScanStatus(userId: string): Observable<void> {
    return new Observable(observer => {
      this.usersService.getScanStatus(userId).subscribe({
        next: (response) => {
          if (response.code === 200) {
            this.scanStatus = response.data;
          }
          observer.next();
          observer.complete();
        },
        error: (err) => {
          console.error('Error fetching scan status', err);
          observer.complete();
        }
      });
    });
  }

  trashStatus: ScanStatus | null = null;
  loadTrashStatus(userId: string): void {
    this.usersService.getTrashStatus(userId).subscribe({
      next: (response) => {
        if (response.code === 200) {
          this.trashStatus = response.data;
        }
      },
      error: (err) => console.error('Error fetching trash status', err)
    });
  }

  toggleErrorSimulation(): void {
    this.loadCurrentUser();
  }

  toggleSelection(result: ScanResult): void {
    if (this.selectedSenders.has(result)) {
      this.selectedSenders.delete(result);
    } else {
      this.selectedSenders.add(result);
    }
  }

  toggleAllSelection(): void {
    if (this.isAllSelected) {
      this.selectedSenders.clear();
    } else {
      this.scanStatus?.Results?.forEach(r => this.selectedSenders.add(r));
    }
  }

  openDeleteModal(): void {
    this.showDeleteModal = true;
  }

  // Trash Job State
  trashJobLoading = false;
  trashJobMessage: string | null = null;
  trashJobError: string | null = null;

  closeDeleteModal(): void {
    this.showDeleteModal = false;
    // Reset trash job state when closing
    this.trashJobLoading = false;
    this.trashJobMessage = null;
    this.trashJobError = null;
  }

  proceedWithDelete(): void {
    if (!this.user) return;

    this.trashJobLoading = true;
    this.trashJobMessage = null;
    this.trashJobError = null;

    const sendersToDelete = this.selectedSendersList.map(s => s.SenderEmail);

    this.jobsService.startTrashJob(this.user.UserId, sendersToDelete).subscribe({
      next: (response) => {
        this.trashJobLoading = false;
        if (response.code === 200) {
          // Force local update to Pending state immediately
          this.trashStatus = {
            IsPending: true,
            IsSuccess: false,
            IsError: false,
            Results: [],
            JobId: 'new-job', // placeholder
            NoExistingJobs: false
          };
          this.closeDeleteModal();
        } else {
          this.trashJobError = response.msg || 'Failed to start deletion job.';
        }
      },
      error: (err) => {
        this.trashJobLoading = false;
        console.error('Error starting trash job', err);
        this.trashJobError = 'Network error while starting deletion job.';
      }
    });
  }

  startScan(): void {
    if (!this.user) return;
    this.scanMessage = null;
    this.scanErrorMessage = null;
    this.simulatedScanSuccess = false;


    this.jobsService.startScan(this.user.UserId).subscribe({
      next: (response) => {
        if (response.code === 200) {
          this.scanMessage = response.data;
          this.simulatedScanSuccess = true; // Use simulated success to trigger UI update for now
          if (this.user) {
            this.loadScanStatus(this.user.UserId);
          }
        }
      },
      error: (err) => {
        console.error('Error starting scan', err);
        this.scanErrorMessage = 'Failed to start scan job.';
        this.simulatedScanError = true;
      }
    });
  }

  toggleScanSuccess(): void {
    if (this.simulatedScanSuccess) {
      this.scanMessage = "Simulated: Started the job to scan your mailbox";
      this.scanErrorMessage = null;
      this.simulatedScanError = false; // Mutually exclusive
    } else {
      this.scanMessage = null;
    }
  }

  toggleScanError(): void {
    if (this.simulatedScanError) {
      this.scanErrorMessage = "Simulated: Failed to start scan job.";
      this.scanMessage = null;
      this.simulatedScanSuccess = false; // Mutually exclusive
    } else {
      this.scanErrorMessage = null;
    }
  }
}
