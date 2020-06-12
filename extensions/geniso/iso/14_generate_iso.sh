#!/bin/sh

set -ex



# Generate ISO image for UEFI based systems.
uefi() {
  cd $ISOIMAGE

  # Now we generate 'hybrid' ISO image file which can also be used on
  # USB flash drive, e.g. 'dd if=minimal_linux_live.iso of=/dev/sdb'.
  xorriso -as mkisofs \
    -isohybrid-mbr $ISOIMAGE/boot/syslinux/isohdpfx.bin \
    -c boot/boot.cat \
    -e boot/uefi.img \
      -no-emul-boot \
      -isohybrid-gpt-basdat \
    -o "$ROOT_DIR/${IMAGE_NAME}" \
    $ISOIMAGE
}

# Generate ISO image for BIOS based systems.
bios() {
  cd $ISOIMAGE

  # Now we generate 'hybrid' ISO image file which can also be used on
  # USB flash drive, e.g. 'dd if=minimal_linux_live.iso of=/dev/sdb'.
  xorriso -as mkisofs \
    -isohybrid-mbr $ISOIMAGE/boot/syslinux/isohdpfx.bin \
    -c boot/syslinux/boot.cat \
    -b boot/syslinux/isolinux.bin \
      -no-emul-boot \
      -boot-load-size 4 \
      -boot-info-table \
    -o "$ROOT_DIR/${IMAGE_NAME}" \
    $ISOIMAGE
}

# Generate ISO image for both BIOS and UEFI based systems.
both() {
  cd $ISOIMAGE

  xorriso -as mkisofs \
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

echo "*** GENERATE ISO BEGIN ***"

if [ ! -d $ISOIMAGE ] ; then
  echo "Cannot locate ISO image work folder. Cannot continue."
  exit 1
fi

case $FIRMWARE_TYPE in
  bios)
    bios
    ;;

  uefi)
    uefi
    ;;

  both)
    both
    ;;

  *)
    echo "Firmware type '$FIRMWARE_TYPE' is not recognized. Cannot continue."
    exit 1
    ;;
esac

cat << CEOF

  #################################################################
  #                                                               #
  #  ISO image file '$IMAGE_NAME.iso' has been generated.  #
  #                                                               #
  #################################################################

CEOF
