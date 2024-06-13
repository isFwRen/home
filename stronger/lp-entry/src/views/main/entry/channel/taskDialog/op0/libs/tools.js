export const getSrc = (bill, fileUrl, index) => {
  const { downloadPath, pictures } = bill
  const path = pictures[index] || pictures[0]

  return `${ fileUrl }${ downloadPath }${ path }`
}