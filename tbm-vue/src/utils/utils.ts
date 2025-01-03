export function timeSinceString(unixMillis: number) {  
    if (!unixMillis) {
      return "never";
    }

    const unixSeconds = unixMillis / 1000;
  
    const now = Math.floor(Date.now() / 1000); // Get current time in seconds
    const seconds = now - unixSeconds; // Calculate the difference
  
    if (seconds < 0) {
      return "just now"; // Handle future dates
    }
  
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
  
    if (days > 0) {
      return days + " days ago";
    } else if (hours > 0) {
      return hours + " hours ago";
    } else if (minutes > 0) {
      return minutes + " minutes ago";
    } else {
      return "just now";
    }
  }

  export function timeSinceSeconds(unixMillis: number) {
    const thenSeconds = unixMillis / 1000;
    const now = Math.floor(Date.now() / 1000); // Get current time in seconds
    const seconds = now - thenSeconds; // Calculate the difference

    return seconds
  }