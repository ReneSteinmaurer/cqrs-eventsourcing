import {Component, viewChild} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {MatButtonModule} from '@angular/material/button';
import {NavbarComponent} from './core/navbar/navbar.component';
import {MatIconModule} from '@angular/material/icon';

@Component({
  selector: 'app-root',
  imports: [MatButtonModule, RouterOutlet, NavbarComponent, MatIconModule],
  templateUrl: './app.component.html',
})
export class AppComponent {
}
