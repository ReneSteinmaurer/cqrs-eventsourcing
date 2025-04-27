import {Component, ElementRef, inject, viewChild} from '@angular/core';
import {AgGridAngular, AgGridModule} from 'ag-grid-angular';
import {
  AllCommunityModule,
  ClientSideRowModelModule,
  ColDef,
  colorSchemeDark,
  GridApi,
  GridReadyEvent,
  ModuleRegistry,
  themeAlpine
} from 'ag-grid-community';
import {LandingPageService} from './landing-page.service';
import {AktionenIconRenderer} from './katalogisiere-icon-renderer/aktionen-icon-renderer.component';
import {MatDialog, MatDialogModule} from '@angular/material/dialog';
import {KatalogisiereMediumDialogComponent} from './katalogisiere-medium-dialog/katalogisiere-medium-dialog.component';
import {MatFormField, MatInput, MatSuffix} from '@angular/material/input';
import {MatIconButton} from '@angular/material/button';
import {MatIcon} from '@angular/material/icon';
import {MatLabel} from '@angular/material/form-field';
import {NgClass} from '@angular/common';
import {Router} from '@angular/router';

ModuleRegistry.registerModules([AllCommunityModule, ClientSideRowModelModule]);
const myTheme = themeAlpine.withPart(colorSchemeDark).withParams({accentColor: "#00fbfb"})

@Component({
  selector: 'app-landing-page',
  imports: [
    AgGridModule,
    AgGridAngular,
    MatDialogModule,
    MatFormField,
    MatInput,
    MatIconButton,
    MatSuffix,
    MatIcon,
    MatLabel,
    NgClass
  ],
  template: `
    <div class="flex bg-[#004f4f] p-4 justify-center items-center">
      <h1 class="flex">Startseite</h1>
    </div>
    <h3 class="flex justify-center items-center mt-8">Bestand</h3>

    <div class="flex ml-4 mr-4 justify-center items-center">
      <div class="w-3/4">
        <div class="flex flex-grow justify-between">
          <div></div>
          <div class="p-2">
            <mat-form-field appearance="outline" class="w-full">
              <mat-label>Suche</mat-label>
              <input
                matInput
                type="text"
                (input)="onFilterTextBoxChanged()"
                #filter
                placeholder="Schnell filtern..."
              />
                <button
                  [ngClass]="{'visible': filter.value}"
                  class="invisible"
                  matSuffix
                  mat-icon-button
                  aria-label="Clear"
                  (click)="resetFilter()"
                >
                  <mat-icon>close</mat-icon>
                </button>
            </mat-form-field>
          </div>
        </div>
        <ag-grid-angular
          class="h-[50vh]"
          domLayout="normal"
          (gridReady)="onGridReady($event)"
          [defaultColDef]="defaultColDef"
          [theme]="myTheme"
          [rowData]="landingPageService.medien.value()"
          [columnDefs]="columnDefs"
          [pagination]="true"
          [paginationPageSize]="20"
        >
        </ag-grid-angular>
      </div>

    </div>

  `,
  styles: ``
})
export class LandingPageComponent {
  router = inject(Router)
  filter = viewChild<ElementRef<HTMLInputElement>>('filter')
  private gridApi!: GridApi
  landingPageService = inject(LandingPageService);
  dialog = inject(MatDialog);

  openKatalogisierenDialog(mediumId: string) {
    this.dialog.open(KatalogisiereMediumDialogComponent, {
      width: '70%',
      height: '70%',
      data: {
        mediumId
      }
    })
  }

  private openDetailansicht(mediumId: string) {
    this.router.navigate(['medium', mediumId])
  }

  columnDefs: ColDef[] = [
    { field: 'mediumId', headerName: 'ID', hide: true },
    { field: 'name', headerName: 'Titel', lockPosition: true, flex: 2, minWidth: 250 },
    { field: 'ISBN', headerName: 'ISBN', lockPosition: true, flex: 1, minWidth: 150 },
    { field: 'mediumType', headerName: 'Typ', flex: 1 },
    { field: 'genre', headerName: 'Genre', flex: 1 },
    { field: 'standort', headerName: 'Bibliothek', flex: 1 },
    {
      headerName: 'Aktionen',
      pinned: 'right',
      lockPosition: true,
      width: 200,
      minWidth: 150,
      maxWidth: 250,
      cellRendererSelector: () => {
        return {
          component: AktionenIconRenderer,
          params: {
            katalogisieren: (mediumId: string) => this.openKatalogisierenDialog(mediumId),
            detailansichtOeffnen: (mediumId: string) => this.openDetailansicht(mediumId),
          },
        };
      },
    },
  ];

  defaultColDef: ColDef = {
    autoHeight: true,
    sortable: true,
    filter: true,
    resizable: false,
  };

  onFilterTextBoxChanged() {
    if (!this.filter()) {
      return
    }
    this.gridApi.setGridOption(
      "quickFilterText",
      this.filter()!.nativeElement.value,
    );
  }

  onGridReady(params: GridReadyEvent<any>) {
    this.gridApi = params.api;
  }

  resetFilter() {
    if (!this.filter()) {
      return
    }
    this.filter()!.nativeElement.value = '';
    this.onFilterTextBoxChanged()
  }

  protected readonly myTheme = myTheme;
}
