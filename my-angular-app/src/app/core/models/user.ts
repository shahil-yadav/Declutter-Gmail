export interface User {
    UserId: string;
    Email: string;
    FullName: string;
    CreatedAt: string;
    CoverPhoto: string;
    access_token: string;
    expiry: string;
}

export interface ApiResponse<T> {
    code: number;
    msg: string;
    data: T;
}

export interface ScanResult {
    SenderEmail: string;
    Count: number;
}

export interface ScanStatus {
    JobId: string;
    IsPending: boolean;
    IsSuccess: boolean;
    IsError: boolean;
    NoExistingJobs: boolean;
    Results?: ScanResult[];
}
