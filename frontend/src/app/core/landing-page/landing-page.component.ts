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
    <div class="flex w-full flex-col items-center">
      <ag-grid-angular
        style="width: 100%; height: 100%;"
        domLayout="autoHeight"
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
  ];

  defaultColDef: ColDef = {
    autoHeight: true,
    suppressAutoSize: false,
    suppressSizeToFit: false,
    sortable: true,
    filter: true,
    resizable: true,
  };

  protected readonly myTheme = myTheme;
}
