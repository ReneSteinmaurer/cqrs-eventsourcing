import {Component, inject} from '@angular/core';
import {AgGridAngular, AgGridModule} from 'ag-grid-angular';
import {
  AllCommunityModule,
  ClientSideRowModelModule,
  ColDef,
  colorSchemeDark,
  ModuleRegistry,
  themeAlpine
} from 'ag-grid-community';
import { LandingPageService } from './landing-page.service';
import {KatalogisiereIconRendererComponent} from './katalogisiere-icon-renderer/katalogisiere-icon-renderer.component';

ModuleRegistry.registerModules([AllCommunityModule, ClientSideRowModelModule]);
const myTheme = themeAlpine.withPart(colorSchemeDark).withParams({accentColor: "#00fbfb"})

@Component({
  selector: 'app-landing-page',
  imports: [
    AgGridModule,
    AgGridAngular
  ],
  template: `
    <div class="flex bg-[#004f4f] p-4 justify-center items-center">
      <h1 class="flex">Startseite</h1>
    </div>
    <h3 class="flex justify-center items-center mt-8">Bestand</h3>
    <div class="flex m-4 justify-center items-center">
      <ag-grid-angular
        class="w-full h-[50vh]"
        domLayout="normal"
        [defaultColDef]="defaultColDef"
        [theme]="myTheme"
        [rowData]="landingPageService.medien()"
        [columnDefs]="columnDefs"
        [pagination]="true"
        [paginationPageSize]="10"
      >
      </ag-grid-angular>
    </div>

  `,
  styles: ``
})
export class LandingPageComponent {
  landingPageService = inject(LandingPageService);

  constructor() {
    this.landingPageService.getAllMedien()
  }

  columnDefs: ColDef[] = [
    {field: 'mediumId', headerName: 'ID', hide: true},
    {field: 'name', flex: 1, headerName: 'Titel'},
    {field: 'ISBN', flex: 1, headerName: 'ISBN'},
    {field: 'mediumType', headerName: 'Typ'},
    {field: 'genre', headerName: 'Genre'},
    {field: 'signature', headerName: 'Signatur'},
    {field: 'standort', headerName: 'Bibliothek'},
    {field: 'exemplarCode', headerName: 'Exemplar-ID'},
    {
      headerName: 'Katalogisieren',
      cellRenderer: KatalogisiereIconRendererComponent,
      pinned: 'right',
      width: 150,
      minWidth: 90,
      maxWidth: 150,
      cellStyle: {
        border: 'none',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
      },
    },
  ];

  defaultColDef: ColDef = {
    autoHeight: true,
    sortable: true,
    filter: true,
    resizable: false,
  };

  protected readonly myTheme = myTheme;
}
