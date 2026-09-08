const BOX = { x: 34, y: 213, width: 337.64, height: 215.28 };
const VIDEO = { streamWidth: 404, streamHeight: 720 };

const extendBoxSize = (boxSize, multiplier, maxWidth, maxHeight) => {
  //TODO: make newWeidth/newHeight no more than stream size
  const newWidth = Math.min(maxWidth, Math.floor(boxSize.width * multiplier));
  const newHeight = Math.min(
    maxHeight,
    Math.floor(boxSize.height * multiplier),
  );
  const newX = Math.max(
    0,
    Math.floor(boxSize.x - (newWidth - boxSize.width) / 2),
  );
  const newY = Math.max(
    0,
    Math.floor(boxSize.y - (newHeight - boxSize.height) / 2),
  );

  return { x: newX, y: newY, width: newWidth, height: newHeight };
};

const normalizeCoordinates = (boxSize, videoSize) => {
  const SIZE_INCREASE_FOR_IMAGE_ENHANCEMENT = 1.2;
  const { x, y, width, height } = extendBoxSize(
    boxSize,
    SIZE_INCREASE_FOR_IMAGE_ENHANCEMENT,
    videoSize.streamWidth,
    videoSize.streamHeight,
  );
  console.log(x, y, width, height);

  const normX = x / videoSize.streamWidth;
  const normY = y / videoSize.streamHeight;
  const normWidth = width / videoSize.streamWidth;
  const normHeight = height / videoSize.streamHeight;

  return [
    [normX, normY],
    [normWidth, normHeight],
  ];
};

console.log(normalizeCoordinates(BOX, VIDEO));
