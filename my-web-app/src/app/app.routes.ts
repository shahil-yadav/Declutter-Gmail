import { Routes } from '@angular/router';
import { HomePage } from './home/home.page';
import { UserPage } from './user/user.page';

export const routes: Routes = [
  {
    path: '',
    component: HomePage,
  },
  {
    path: ':id',
    component: UserPage,
  },
];
