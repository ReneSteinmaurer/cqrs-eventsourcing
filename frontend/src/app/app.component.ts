import {Component} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {MatButtonModule} from '@angular/material/button';

@Component({
  selector: 'app-root',
  imports: [MatButtonModule, RouterOutlet],
  templateUrl: './app.component.html',
})
export class AppComponent {
  tasks = [
    { label: 'Nutzer registrieren', route: '/registrieren', icon: 'person_add' },
    { label: 'Medium erwerben', route: '/medium-erwerben', icon: 'library_add' },
    { label: 'Medium katalogisieren', route: '/medium-katalogisieren', icon: 'search' },
    { label: 'Medium verleihen', route: '/medium-verleihen', icon: 'assignment_return' },
    { label: 'Medium zurückgeben', route: '/medium-zurueckgeben', icon: 'assignment_turned_in' },
    { label: 'Medienübersicht', route: '/medien', icon: 'view_list' },
  ];

  navigate(route: string) {
    // Router navigation hier, z.B. this.router.navigateByUrl(route)
  }
}
