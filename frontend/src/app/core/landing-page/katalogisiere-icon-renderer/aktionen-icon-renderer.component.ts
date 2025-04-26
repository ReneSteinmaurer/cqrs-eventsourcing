import {Component, signal} from '@angular/core';
import {ICellRendererAngularComp} from 'ag-grid-angular';
import {ICellRendererParams} from 'ag-grid-community';
import {MatIconButton} from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';
import {MediumBestand} from '../../../shared/types/medium-bestand';
import { MatTooltip } from '@angular/material/tooltip';

export interface AktionenCellRendererParams extends ICellRendererParams<MediumBestand> {
  katalogisieren: (mediumId: string) => void;
  verleihen: (mediumId: string) => void;
  detailansichtOeffnen: (mediumId: string) => void;
}

@Component({
  selector: 'app-katalogisiere-icon-renderer',
  imports: [
    MatIconModule,
    MatIconButton,
    MatTooltip
  ],
  template: `
    @if (!medium()?.katalogisiert) {
      <button
        class="flex text-primary justify-center items-center"
        mat-icon-button
        matTooltip="Medium katalogisieren"
        (click)="katalogisieren(medium()!.mediumId)"
      >
        <mat-icon>library_add</mat-icon>
      </button>
    }
    @if (medium()?.katalogisiert) {
      <button
        class="flex text-primary justify-center items-center"
        mat-icon-button
        matTooltip="Medium verleihen"
        (click)="verleihen(medium()!.mediumId)"
      >
        <mat-icon>label_important</mat-icon>
      </button>
    }
    <button
      class="flex text-primary justify-center items-center"
      mat-icon-button
      matTooltip="Details ansehen"
      (click)="detailansichtOeffnen(medium()!.mediumId)"
    >
      <mat-icon>info</mat-icon>
    </button>

  `,
  styles: `
    .mat-mdc-icon-button {
      font-size: 0 !important;
    }
  `
})
export class AktionenIconRenderer implements ICellRendererAngularComp {
  medium = signal<MediumBestand | undefined>(undefined)
  katalogisieren = (id: string) => {};
  verleihen = (id: string) => {};
  detailansichtOeffnen = (id: string) => {};

  agInit(params: AktionenCellRendererParams): void {
    this.medium.set(params.data)
    this.katalogisieren = params.katalogisieren;
    this.verleihen = params.verleihen;
    this.detailansichtOeffnen = params.detailansichtOeffnen;
    this.refresh(params);
  }

  refresh(params: ICellRendererParams): boolean {
    return true;
  }
}
