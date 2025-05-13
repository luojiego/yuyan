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
    $('body').append('<div id="loading-overlay" style="position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.5); z-index: 9999; display: flex; justify-content: center; align-items: center;"><div class="spinner-border text-light" role="status"><span class="sr-only">' + __('loading') + '</span></div></div>');
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
        <strong class="mr-auto">${__('notification')}</strong>
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
  
  let errorMessage = __('error_occurred');
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
  showToast(__('copied_clipboard'));
}

// Add event handlers after page load
$(document).ready(function() {
  // Apply language settings
  updatePageLanguage();
  
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
    if (!confirm(__('confirm_delete'))) {
      e.preventDefault();
      e.stopPropagation();
      return false;
    }
  });
  
  // Set active navigation item
  setActiveNavItem();
  
  // Load dashboard data if on dashboard page
  if (window.location.pathname === '/') {
    loadDashboardData();
  }
  
  // Set up quick send form submission on dashboard
  $('#quick-send-form').on('submit', function(e) {
    e.preventDefault();
    sendQuickMessage();
  });
  
  // Handle language switching
  $('.language-option').on('click', function(e) {
    e.preventDefault();
    const lang = $(this).data('lang');
    if (changeLanguage(lang)) {
      // Language changed successfully
      console.log('Language changed to ' + lang);
    }
  });
  
  // Listen for language change events
  document.addEventListener('languageChanged', function(e) {
    console.log('Language changed event received: ' + e.detail.language);
    // Update any dynamic content that needs language updates
    if (window.location.pathname === '/') {
      loadDashboardData();
    } else if (window.location.pathname.includes('/bots')) {
      loadBots();
    } else if (window.location.pathname.includes('/messages')) {
      loadMessages();
    }
  });
});

// Dashboard functions
function loadDashboardData() {
  // Load bot data
  $.ajax({
    url: '/api/bots',
    method: 'GET',
    success: function(bots) {
      $('#total-bots').text(bots.length || 0);
      
      // Populate bot select in quick send form
      const botSelect = $('#bot-select');
      if (botSelect.length) {
        botSelect.empty();
        if (bots.length === 0) {
          botSelect.append(`<option value="">${__('no_bots_available')}</option>`);
        } else {
          bots.forEach(function(bot) {
            if (bot.is_active) {
              botSelect.append(`<option value="${bot.id}">${bot.name} (${__(bot.type)})</option>`);
            }
          });
        }
      }
    },
    error: function(xhr) {
      handleAjaxError(xhr);
    }
  });
  
  // Load message data for stats and recent messages
  $.ajax({
    url: '/api/messages',
    method: 'GET',
    success: function(messages) {
      // Calculate stats
      let sentCount = 0;
      let pendingCount = 0;
      let failedCount = 0;
      
      if (messages && messages.length > 0) {
        messages.forEach(function(message) {
          if (message.status === 'sent') {
            sentCount++;
          } else if (message.status === 'pending' || message.status === 'processing') {
            pendingCount++;
          } else if (message.status === 'failed') {
            failedCount++;
          }
        });
      }
      
      // Update stats display
      $('#total-messages-sent').text(sentCount);
      $('#total-messages-pending').text(pendingCount);
      $('#total-messages-failed').text(failedCount);
      
      // Populate recent messages table
      const tbody = $('#recent-messages');
      if (tbody.length) {
        tbody.empty();
        
        if (messages && messages.length > 0) {
          // Get 5 most recent messages
          const recentMessages = messages.slice(0, 5);
          
          recentMessages.forEach(function(message) {
            let statusClass = 'secondary';
            if (message.status === 'sent') statusClass = 'success';
            if (message.status === 'failed') statusClass = 'danger';
            if (message.status === 'pending') statusClass = 'warning';
            
            // Get bot name, using message.bot.name if available, otherwise try message.bot_name or 'Unknown'
            let botName = __('unknown');
            if (message.bot && message.bot.name) {
              botName = message.bot.name;
            } else if (message.bot_name) {
              botName = message.bot_name;
            }
            
            const row = `
              <tr>
                <td>${botName}</td>
                <td><span class="badge badge-${statusClass}">${__(message.status)}</span></td>
                <td>${new Date(message.sent_at).toLocaleString()}</td>
                <td>
                  <a href="/messages?id=${message.id}" class="btn btn-sm btn-info">
                    <i class="fas fa-eye"></i>
                  </a>
                </td>
              </tr>
            `;
            tbody.append(row);
          });
        } else {
          tbody.append(`<tr><td colspan="4" class="text-center">${__('no_messages_found')}</td></tr>`);
        }
      }
      
      // Apply translations to any newly added elements
      updatePageLanguage();
    },
    error: function(xhr) {
      handleAjaxError(xhr);
    }
  });
}

// Set active navigation item based on current path
function setActiveNavItem() {
  const path = window.location.pathname;
  if (path === '/') {
    $('#nav-dashboard').addClass('active');
  } else if (path.includes('/bots')) {
    $('#nav-bots').addClass('active');
  } else if (path.includes('/messages')) {
    $('#nav-messages').addClass('active');
  } else if (path.includes('/settings')) {
    $('#nav-settings').addClass('active');
  }
}

// Handle quick send message form submission
function sendQuickMessage() {
  showLoading();
  
  const botId = $('#bot-select').val();
  const format = $('#message-format').val();
  const content = $('#message-content').val();
  
  if (!botId || !content) {
    hideLoading();
    showToast(__('please_select_bot_enter_message'), 'danger');
    return;
  }
  
  $.ajax({
    url: '/api/messages',
    method: 'POST',
    contentType: 'application/json',
    data: JSON.stringify({
      bot_id: parseInt(botId),
      content: content,
      format: format
    }),
    success: function(response) {
      hideLoading();
      showToast(__('message_sent_successfully'), 'success');
      $('#message-content').val(''); // Clear the message content
      loadDashboardData(); // Reload dashboard data to show the new message
    },
    error: function(xhr) {
      handleAjaxError(xhr);
    }
  });
} 