import { Component, inject } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { ScanJobResults, Status, TrashJobResults, User } from '../types';
import { UsersService } from '../users.service';
import { JobsService } from '../jobs.service';

@Component({
  selector: 'app-user',
  imports: [RouterLink],
  templateUrl: './user.page.html',
  styleUrl: './user.page.css',
})
export class UserPage {
  private activatedRoute = inject(ActivatedRoute);
  private userService = inject(UsersService);
  private jobsService = inject(JobsService);

  senders: string[] = [];
  scanJobInfo?: ScanJobResults;
  trashJobInfo?: TrashJobResults;
  user?: User;
  userId?: string;

  bindSender(sender: string) {
    if (this.senders.includes(sender)) {
      this.senders = this.senders.filter((v) => v !== sender);
    } else {
      this.senders.push(sender);
    }
  }

  async createTrashJob() {
    if (this.userId) {
      const resp = await this.jobsService.createTrashJob(this.userId, this.senders);
      alert(resp.data);
      window.location.reload();
    } else {
      console.error('User-id not set');
    }
  }

  async createScanJob() {
    if (this.userId) {
      const resp = await this.jobsService.createScanJob(this.userId);
      alert(resp.data);

      window.location.reload();
    }
  }

  constructor() {
    this.activatedRoute.params.subscribe(async (params) => {
      const userId = params['id'] as string;
      // storing the id in a instance variable
      this.userId = userId;

      // setting user variable
      const userResp = await this.userService.getUser(userId);
      this.user = userResp.data;

      // setting scan jobs results
      const scanJobInfoResp = await this.jobsService.getScanJobInfo(userId);
      this.scanJobInfo = scanJobInfoResp.data;

      // setting trash jobs results
      const trashJobInfoResp = await this.jobsService.getTrashJobInfo(userId);
      this.trashJobInfo = trashJobInfoResp.data;
    });
  }
}
