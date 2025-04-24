import {Component, signal} from '@angular/core';
import {ICellRendererAngularComp} from 'ag-grid-angular';
import {ICellRendererParams} from 'ag-grid-community';
import {MatIconButton} from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';

@Component({
  selector: 'app-katalogisiere-icon-renderer',
  imports: [
    MatIconModule,
    MatIconButton
  ],
  template: `
    <button class="flex justify-center items-center" mat-icon-button>
      <mat-icon>bolt</mat-icon>
    </button>
  `,
  styles: ``
})
export class KatalogisiereIconRendererComponent implements ICellRendererAngularComp {

  agInit(params: ICellRendererParams<any, any, any>): void {
    this.refresh(params);
  }

  refresh(params: ICellRendererParams<any, any, any>): boolean {
    return true;
  }

}
