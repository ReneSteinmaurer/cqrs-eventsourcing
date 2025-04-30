import {Component, inject} from '@angular/core';
import {
  ConfirmationDialogComponent
} from '../../../../../../shared/ui/confirmation-dialog/confirmation-dialog.component';
import {MatDialogRef} from '@angular/material/dialog';
import {MediumDetailService} from '../../../services/medium-detail.service';

@Component({
  selector: 'app-wiederaufgefunden-durch-nutzer-dialog',
  imports: [
    ConfirmationDialogComponent
  ],
  template: `
    <confirmation-dialog (close)="close()" (confirm)="confirm()" titel="Medium Bestandsverlust">
      <p #content>
        Möchtest du dieses Medium wirklich als wiederaufgefunden markieren?
      </p>
    </confirmation-dialog>
  `,
})
export class WiederaufgefundenDurchNutzerDialogComponent {
  private dialogRef = inject(MatDialogRef<WiederaufgefundenDurchNutzerDialogComponent>);
  detailService = inject(MediumDetailService)

  close() {
    this.dialogRef.close()
  }

  confirm() {
    this.detailService.wiederaufgefundenDurchNutzer().subscribe(() => {
      this.dialogRef.close()
    });
  }
}
