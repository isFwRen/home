import QRCode from 'qrcode'

export const getQRcode = (qrCodeStr, width) => {
  let imgUrl= ''
    QRCode.toDataURL(
        qrCodeStr,
      { errorCorrectionLevel: 'L', margin: 2, width },
      (err, url) => {
        if (err) throw err;
        imgUrl =  url
      }
    )
    return imgUrl
  }
  