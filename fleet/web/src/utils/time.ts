export const getRelativeTimestamp = (timestamp) => {
  const now = new Date();
  const then = new Date(timestamp);

  const diff = now - then;

  if (diff < 60000) {
    return 'just now';
  }

  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);

  if (days > 0) {
    return `${days}days ago`;
  } else if (hours > 0) {
    return `${hours}hours ago`;
  } else if (minutes > 0) {
    return `${minutes}minutes ago`;
  } else {
    return 'just now';
  }
};
