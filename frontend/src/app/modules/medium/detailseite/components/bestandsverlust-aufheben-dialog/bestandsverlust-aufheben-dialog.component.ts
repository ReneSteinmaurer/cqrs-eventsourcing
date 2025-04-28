import {Component, inject} from '@angular/core';
import {ConfirmationDialogComponent} from '../../../../../shared/ui/confirmation-dialog/confirmation-dialog.component';
import {MatDialogRef} from '@angular/material/dialog';
import {MediumDetailService} from '../../services/medium-detail.service';

@Component({
  selector: 'app-bestandsverlust-aufheben-dialog',
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
  styles: ``
})
export class BestandsverlustAufhebenDialogComponent {
  private dialogRef = inject(MatDialogRef<BestandsverlustAufhebenDialogComponent>);
  detailService = inject(MediumDetailService)

  close() {
    this.dialogRef.close()
  }

  confirm() {
    this.detailService.bestandsverlustAufheben().subscribe(() => {
      this.dialogRef.close()
    });
  }
}
