import { Routes } from '@angular/router';
import { DashboardComponent } from './features/dashboard/dashboard';
import { UserDetailsComponent } from './features/user-details/user-details';

export const routes: Routes = [
    { path: '', component: DashboardComponent },
    { path: 'users/:id', component: UserDetailsComponent },
];
