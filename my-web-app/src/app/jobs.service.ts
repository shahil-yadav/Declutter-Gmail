import { Injectable } from '@angular/core';
import { Response, ScanJobResults, TrashJobResults } from './types';
import axios from 'axios';

type ScanJobResultsResponse = Response<ScanJobResults>;
type TrashJobResultsResponse = Response<TrashJobResults>;
type PostScanJobResponse = Response<string>;
type PostTrashJobResponse = Response<string>;

@Injectable({
  providedIn: 'root',
})
export class JobsService {
  private url = 'http://127.0.0.1:7331/v1';

  async getScanJobInfo(userId: string): Promise<ScanJobResultsResponse> {
    const resp = await axios.get<ScanJobResultsResponse>(`${this.url}/users/${userId}/info/scan`);
    return resp.data;
  }

  async createTrashJob(userId: string, senders: string[]) {
    const urlencoded = new URLSearchParams();
    for (const sender of senders) {
      urlencoded.append('sender[]', sender);
    }
    urlencoded.append('user-id', userId);

    const resp = await axios.post<PostTrashJobResponse>(
      `${this.url}/job/trash`,
      urlencoded.toString(),
      {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      },
    );

    return resp.data;
  }

  async createScanJob(userId: string): Promise<PostScanJobResponse> {
    const resp = await axios.post<PostScanJobResponse>(
      `${this.url}/job/scan`,
      { 'user-id': userId },
      {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      },
    );

    return resp.data;
  }

  async getTrashJobInfo(userId: string): Promise<TrashJobResultsResponse> {
    const resp = await axios.get<TrashJobResultsResponse>(`${this.url}/users/${userId}/info/trash`);
    return resp.data;
  }
}
