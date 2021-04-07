export const findMissingElement = (largeArray, smallArray) => {
  if (smallArray.length === 0) {
    return largeArray[0];
  }
  for (let i = 0; i < largeArray.length; i++) {
    let currentIndicator = largeArray[i];
    for (let j = 0; j < smallArray.length; j++) {
      if (currentIndicator === smallArray[j]) {
        break;
      } else if (j === smallArray.length - 1) {
        return currentIndicator;
      }
    }
  }
};
