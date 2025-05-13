/**
 * Common JavaScript functions for the notification service
 */

// Format date
function formatDate(dateString) {
  if (!dateString) return 'N/A';
  return new Date(dateString).toLocaleString();
}

// Truncate text with ellipsis
function truncateText(text, maxLength) {
  if (!text) return '';
  if (text.length <= maxLength) return text;
  return text.substring(0, maxLength) + '...';
}

// Show loading overlay
function showLoading() {
  if (!$('#loading-overlay').length) {
    $('body').append('<div id="loading-overlay" style="position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.5); z-index: 9999; display: flex; justify-content: center; align-items: center;"><div class="spinner-border text-light" role="status"><span class="sr-only">Loading...</span></div></div>');
  } else {
    $('#loading-overlay').show();
  }
}

// Hide loading overlay
function hideLoading() {
  $('#loading-overlay').hide();
}

// Show toast notification
function showToast(message, type = 'success') {
  // Remove existing toasts
  $('.toast').remove();
  
  // Create new toast
  const toast = $(`
    <div class="toast" role="alert" aria-live="assertive" aria-atomic="true" style="position: fixed; top: 20px; right: 20px; z-index: 9999;">
      <div class="toast-header bg-${type} text-white">
        <strong class="mr-auto">Notification</strong>
        <button type="button" class="ml-2 mb-1 close" data-dismiss="toast" aria-label="Close">
          <span aria-hidden="true">&times;</span>
        </button>
      </div>
      <div class="toast-body">
        ${message}
      </div>
    </div>
  `);
  
  // Add to body and show
  $('body').append(toast);
  toast.toast({
    delay: 3000,
    autohide: true
  });
  toast.toast('show');
}

// Handle AJAX errors
function handleAjaxError(xhr) {
  hideLoading();
  
  let errorMessage = 'An error occurred';
  if (xhr.responseJSON && xhr.responseJSON.error) {
    errorMessage = xhr.responseJSON.error;
  } else if (xhr.statusText) {
    errorMessage = xhr.statusText;
  }
  
  showToast(errorMessage, 'danger');
}

// Copy text to clipboard
function copyToClipboard(text) {
  const textarea = document.createElement('textarea');
  textarea.value = text;
  document.body.appendChild(textarea);
  textarea.select();
  document.execCommand('copy');
  document.body.removeChild(textarea);
  showToast('Copied to clipboard!');
}

// Add event handlers after page load
$(document).ready(function() {
  // Add copy button to pre elements
  $('pre').each(function() {
    const pre = $(this);
    const copyBtn = $('<button class="btn btn-sm btn-outline-secondary copy-btn" style="position: absolute; top: 5px; right: 5px;"><i class="fas fa-copy"></i></button>');
    
    pre.css('position', 'relative');
    pre.append(copyBtn);
    
    copyBtn.on('click', function() {
      copyToClipboard(pre.text());
    });
  });
  
  // Add confirm dialog to delete buttons
  $(document).on('click', '.btn-delete', function(e) {
    if (!confirm('Are you sure you want to delete this item?')) {
      e.preventDefault();
      e.stopPropagation();
      return false;
    }
  });
}); 