#!/bin/bash

set -e

DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/func.sh"

# Generate ISO image for both BIOS and UEFI based systems.
geniso() {
  cd $ISOIMAGE

  xorriso -as mkisofs \
    -volid "$ISOLABEL" \
    -isohybrid-mbr $ISOIMAGE/boot/syslinux/isohdpfx.bin \
    -c boot/syslinux/boot.cat \
    -b boot/syslinux/isolinux.bin \
      -no-emul-boot \
      -boot-load-size 4 \
      -boot-info-table \
    -eltorito-alt-boot \
    -e boot/uefi.img \
      -no-emul-boot \
      -isohybrid-gpt-basdat \
    -o "$ROOT_DIR/${IMAGE_NAME}" \
  $ISOIMAGE
}

info "Starting generating ISO"

if [ ! -d $ISOIMAGE ] ; then
  echo "Cannot locate ISO image work folder. Cannot continue."
  exit 1
fi

geniso

ok "$IMAGE_NAME has been generated"
