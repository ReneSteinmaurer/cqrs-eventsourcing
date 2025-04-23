import {Component} from '@angular/core';
import {MatDrawer, MatDrawerContainer, MatDrawerContent} from '@angular/material/sidenav';
import {MatButton} from '@angular/material/button';

@Component({
  selector: 'app-navbar',
  template: `
    <mat-drawer-container class="w-full h-full" autosize>
      <mat-drawer #drawer mode="side">
        <p>Auto-resizing sidenav</p>
        @if (showFiller) {
          <p>Lorem, ipsum dolor sit amet consectetur.</p>
        }
        <button (click)="showFiller = !showFiller" mat-raised-button>
          Toggle extra text
        </button>
      </mat-drawer>
      <div class="flex">
        <button type="button" mat-button (click)="drawer.toggle()">
          Toggle sidenav
        </button>
      </div>
    </mat-drawer-container>
    <mat-drawer-content>
      <ng-content></ng-content>
    </mat-drawer-content>
  `,
  imports: [
    MatDrawerContainer,
    MatButton,
    MatDrawer,
    MatDrawerContent
  ]
})
export class NavbarComponent {
  showFiller = false;

}
