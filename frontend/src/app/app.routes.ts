import {Routes} from '@angular/router';
import {LandingPageComponent} from './core/landing-page/landing-page.component';

export const routes: Routes = [
  {
    path: '',
    component: LandingPageComponent
  },
  {
    path: 'medium/:id',
    loadComponent: () => import('./modules/medium/detailseite/detailseite.component')
      .then((x) => x.DetailseiteComponent)
  }
];
