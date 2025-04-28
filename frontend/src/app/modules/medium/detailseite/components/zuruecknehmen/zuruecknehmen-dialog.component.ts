import {Component, inject} from '@angular/core';
import {MatButton} from '@angular/material/button';
import {MatDialogRef} from '@angular/material/dialog';
import {MediumDetailService} from '../../services/medium-detail.service';

@Component({
  selector: 'app-zuruecknehmen',
  imports: [
    MatButton
  ],
  template: `
    <div class="p-6 max-w-full max-h-full">
      <h3 class="text-xl font-semibold m-4">Medium zurückgegeben</h3>
      <p class="text-sm text-gray-400 mt-4 text-center">
        Bist du sicher, dass du das Medium als zurückgenommen markieren möchtest?
      </p>

      <div class="flex justify-end gap-4 w-full pt-4">
        <button mat-button (click)="cancel()">Abbrechen</button>
        <button mat-flat-button color="warn" (click)="confirm()">Bestätigen</button>
      </div>
    </div>
  `,
  styles: ``
})
export class ZuruecknehmenDialogComponent {
  private dialogRef = inject(MatDialogRef<ZuruecknehmenDialogComponent>);
  detailService = inject(MediumDetailService)

  cancel() {
    this.dialogRef.close()
  }

  confirm() {
    this.detailService.zurueckgeben().subscribe(() => {
      this.dialogRef.close()
    })
  }
}
