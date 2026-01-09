export interface Response<T = any> {
  code: number;
  msg: string;
  data: T;
}

export interface User {
  UserId: number;
  Email: string;
  FullName: string;
  CreatedAt: string;
  CoverPhoto: string;
}

export enum Status {
  pending,
  success,
  failure,
}

export interface MailFolder {
  Count: number;
  SenderEmail: string;
}

export interface JobStatus {
  JobId: string;
  IsPending: boolean;
  IsSuccess: boolean;
  IsError: boolean;
}

export interface ScanJobResults extends JobStatus {
  Results: MailFolder[];
}

export interface TrashJobResults extends JobStatus {}
