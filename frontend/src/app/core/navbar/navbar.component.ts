import {Component, viewChild} from '@angular/core';
import {MatDrawer, MatDrawerContainer, MatDrawerContent} from '@angular/material/sidenav';
import {MatButton} from '@angular/material/button';

@Component({
  selector: 'app-navbar',
  template: `
    <mat-drawer-container  autosize>
      <mat-drawer #drawer mode="side" class="flex flex-row">
        <div class="flex justify-center items-center mt-4">
          <button mat-flat-button>Verleihen</button>
        </div>
        <div class="flex justify-center items-center mt-4">
          <button mat-flat-button>Katalogisiere</button>
        </div>
      </mat-drawer>
      <div class="flex">
        <button type="button" mat-button (click)="drawer.toggle()">
          Toggle sidenav
        </button>
      </div>
      <mat-drawer-content>
        <ng-content></ng-content>
      </mat-drawer-content>
    </mat-drawer-container>
  `,
  imports: [
    MatDrawerContainer,
    MatButton,
    MatDrawer,
    MatDrawerContent,
  ]
})
export class NavbarComponent {
  drawer = viewChild.required<MatDrawer>('drawer')
}
