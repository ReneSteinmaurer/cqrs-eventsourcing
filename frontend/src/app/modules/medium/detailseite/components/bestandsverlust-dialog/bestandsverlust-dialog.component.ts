import {Component, inject} from '@angular/core';
import {ConfirmationDialogComponent} from '../../../../../shared/ui/confirmation-dialog/confirmation-dialog.component';
import {MatDialogRef} from '@angular/material/dialog';
import {MediumDetailService} from '../../services/medium-detail.service';

@Component({
  selector: 'app-bestandsverlust-dialog',
  imports: [
    ConfirmationDialogComponent
  ],
  template: `
    <confirmation-dialog (close)="close()" (confirm)="confirm()" titel="Medium Bestandsverlust">
      <p #content>
        Möchtest du dieses Medium wirklich als verloren melden?
      </p>
    </confirmation-dialog>
  `,
  styles: ``
})
export class BestandsverlustDialogComponent {
  private dialogRef = inject(MatDialogRef<BestandsverlustDialogComponent>);
  detailService = inject(MediumDetailService)

  close() {
    this.dialogRef.close()
  }

  confirm() {
    this.detailService.bestandsverlust().subscribe(() => {
      this.dialogRef.close()
    });
  }
}
