/**
 * Internationalization support for the notification service
 */

// Available languages
const availableLanguages = {
  'en': 'English',
  'zh': '中文'
};

// Default language
let currentLanguage = localStorage.getItem('language') || 'en';

// Translations
const translations = {
  'en': {
    // Common
    'notification_service': 'Notification Service',
    'version': 'Version',
    'loading': 'Loading...',
    'confirm_delete': 'Are you sure you want to delete this item?',
    'copied_clipboard': 'Copied to clipboard!',
    'notification': 'Notification',
    'close': 'Close',
    'save': 'Save',
    'update': 'Update',
    'delete': 'Delete',
    'edit': 'Edit',
    'active': 'Active',
    'inactive': 'Inactive',
    'yes': 'Yes',
    'no': 'No',
    'error_occurred': 'An error occurred',
    'unknown': 'Unknown',
    'actions': 'Actions',
    'more_info': 'More info',
    'type': 'Type',
    'format': 'Format',
    'mentions': 'Mentions',
    'created_at': 'Created At',
    'users': 'users',
    
    // Navigation
    'dashboard': 'Dashboard',
    'bot_management': 'Bot Management',
    'message_history': 'Message History',
    'system_settings': 'System Settings',
    
    // Dashboard
    'quick_send': 'Quick Send',
    'total_bots': 'Total Bots',
    'messages_sent': 'Messages Sent',
    'messages_pending': 'Messages Pending',
    'messages_failed': 'Messages Failed',
    'recent_messages': 'Recent Messages',
    'no_messages_found': 'No messages found',
    'no_bots_available': 'No bots available',
    'select_bot': 'Select Bot',
    'message_format': 'Message Format',
    'message_content': 'Message Content',
    'send_message': 'Send Message',
    'please_select_bot_enter_message': 'Please select a bot and enter a message',
    'message_sent_successfully': 'Message sent successfully!',
    
    // Bot Management
    'bot_list': 'Bot List',
    'add_bot': 'Add Bot',
    'add_new_bot': 'Add New Bot',
    'edit_bot': 'Edit Bot',
    'bot_name': 'Bot Name',
    'bot_type': 'Bot Type',
    'bot_token': 'Bot Token',
    'bot_secret': 'Bot Secret',
    'webhook_url': 'Webhook URL',
    'no_bots_found': 'No bots found',
    'confirm_delete_bot': 'Are you sure you want to delete this bot?',
    'failed_load_bots': 'Failed to load bots',
    'failed_create_bot': 'Failed to create bot',
    'failed_update_bot': 'Failed to update bot',
    'failed_delete_bot': 'Failed to delete bot',
    'failed_load_bot_details': 'Failed to load bot details',
    'telegram': 'Telegram',
    'discord': 'Discord',
    'slack': 'Slack',
    'dingtalk': 'DingTalk',
    'custom': 'Custom',
    'dingtalk_secret_help': 'For DingTalk, this is the secret key used for signing webhook requests.',
    'telegram_webhook_help': 'For Telegram, this is the chat ID where messages will be sent.',
    'discord_webhook_help': 'For Discord, enter the complete webhook URL.',
    'slack_webhook_help': 'For Slack, enter the incoming webhook URL.',
    'dingtalk_webhook_help': 'For DingTalk, enter the complete webhook URL from the DingTalk bot settings.',
    'custom_webhook_help': 'Enter the webhook URL for receiving notifications.',
    'dingtalk_token_placeholder': 'Enter access token (optional)',
    'telegram_webhook_placeholder': 'Enter chat ID',
    
    // Message History
    'message_list': 'Message List',
    'sent_at': 'Sent At',
    'status': 'Status',
    'bot': 'Bot',
    'message': 'Message',
    'content': 'Content',
    'view': 'View',
    'message_details': 'Message Details',
    'message_id': 'Message ID',
    'message_status': 'Status',
    'message_sent_at': 'Sent At',
    'retry': 'Retry',
    'sent': 'Sent',
    'pending': 'Pending',
    'processing': 'Processing',
    'failed': 'Failed',
    
    // Settings
    'general_settings': 'General Settings',
    'notification_settings': 'Notification Settings',
    'interface_settings': 'Interface Settings',
    'language': 'Language',
    'theme': 'Theme',
    'dark_mode': 'Dark Mode',
    'light_mode': 'Light Mode',
    'settings_saved': 'Settings saved successfully!',
    'warning': 'Warning',
    'settings_restart_warning': 'Changing these settings will require a server restart to take effect.',
    
    // Server Settings
    'server_settings': 'Server Settings',
    'server_port': 'Server Port',
    'server_mode': 'Server Mode',
    'debug_mode': 'Debug',
    'release_mode': 'Release',
    
    // Database Settings
    'database_settings': 'Database Settings',
    'database_type': 'Database Type',
    'database_file_path': 'Database File Path',
    'database_name': 'Database Name',
    'host': 'Host',
    'port': 'Port',
    'username': 'Username',
    'password': 'Password',
    'connection_parameters': 'Connection Parameters',
    'ssl_mode': 'SSL Mode',
    'disable': 'Disable',
    'require': 'Require',
    'verify_ca': 'Verify CA',
    'verify_full': 'Verify Full'
  },
  'zh': {
    // Common
    'notification_service': '通知服务',
    'version': '版本',
    'loading': '加载中...',
    'confirm_delete': '确定要删除此项目吗？',
    'copied_clipboard': '已复制到剪贴板！',
    'notification': '通知',
    'close': '关闭',
    'save': '保存',
    'update': '更新',
    'delete': '删除',
    'edit': '编辑',
    'active': '活动',
    'inactive': '未活动',
    'yes': '是',
    'no': '否',
    'error_occurred': '发生错误',
    'unknown': '未知',
    'actions': '操作',
    'more_info': '更多信息',
    'type': '类型',
    'format': '格式',
    'mentions': '提及',
    'created_at': '创建时间',
    'users': '用户',
    
    // Navigation
    'dashboard': '仪表盘',
    'bot_management': '机器人管理',
    'message_history': '消息记录',
    'system_settings': '系统设置',
    
    // Dashboard
    'quick_send': '快速发送',
    'total_bots': '机器人总数',
    'messages_sent': '已发送消息',
    'messages_pending': '待处理消息',
    'messages_failed': '失败消息',
    'recent_messages': '最近消息',
    'no_messages_found': '未找到消息',
    'no_bots_available': '没有可用的机器人',
    'select_bot': '选择机器人',
    'message_format': '消息格式',
    'message_content': '消息内容',
    'send_message': '发送消息',
    'please_select_bot_enter_message': '请选择机器人并输入消息',
    'message_sent_successfully': '消息发送成功！',
    
    // Bot Management
    'bot_list': '机器人列表',
    'add_bot': '添加机器人',
    'add_new_bot': '添加新机器人',
    'edit_bot': '编辑机器人',
    'bot_name': '机器人名称',
    'bot_type': '机器人类型',
    'bot_token': '机器人令牌',
    'bot_secret': '机器人密钥',
    'webhook_url': 'Webhook网址',
    'no_bots_found': '未找到机器人',
    'confirm_delete_bot': '确定要删除此机器人吗？',
    'failed_load_bots': '加载机器人失败',
    'failed_create_bot': '创建机器人失败',
    'failed_update_bot': '更新机器人失败',
    'failed_delete_bot': '删除机器人失败',
    'failed_load_bot_details': '加载机器人详情失败',
    'telegram': 'Telegram',
    'discord': 'Discord',
    'slack': 'Slack',
    'dingtalk': '钉钉',
    'custom': '自定义',
    'dingtalk_secret_help': '对于钉钉，这是用于签署webhook请求的密钥。',
    'telegram_webhook_help': '对于Telegram，这是将发送消息的聊天ID。',
    'discord_webhook_help': '对于Discord，请输入完整的webhook URL。',
    'slack_webhook_help': '对于Slack，请输入incoming webhook URL。',
    'dingtalk_webhook_help': '对于钉钉，请输入钉钉机器人设置中的完整webhook URL。',
    'custom_webhook_help': '输入用于接收通知的webhook URL。',
    'dingtalk_token_placeholder': '输入访问令牌（可选）',
    'telegram_webhook_placeholder': '输入聊天ID',
    
    // Message History
    'message_list': '消息列表',
    'sent_at': '发送时间',
    'status': '状态',
    'bot': '机器人',
    'message': '消息',
    'content': '内容',
    'view': '查看',
    'message_details': '消息详情',
    'message_id': '消息ID',
    'message_status': '状态',
    'message_sent_at': '发送时间',
    'retry': '重试',
    'sent': '已发送',
    'pending': '待处理',
    'processing': '处理中',
    'failed': '失败',
    
    // Settings
    'general_settings': '常规设置',
    'notification_settings': '通知设置',
    'interface_settings': '界面设置',
    'language': '语言',
    'theme': '主题',
    'dark_mode': '暗黑模式',
    'light_mode': '明亮模式',
    'settings_saved': '设置保存成功！',
    'warning': '警告',
    'settings_restart_warning': '更改这些设置需要重启服务器才能生效。',
    
    // Server Settings
    'server_settings': '服务器设置',
    'server_port': '服务器端口',
    'server_mode': '服务器模式',
    'debug_mode': '调试模式',
    'release_mode': '发布模式',
    
    // Database Settings
    'database_settings': '数据库设置',
    'database_type': '数据库类型',
    'database_file_path': '数据库文件路径',
    'database_name': '数据库名称',
    'host': '主机',
    'port': '端口',
    'username': '用户名',
    'password': '密码',
    'connection_parameters': '连接参数',
    'ssl_mode': 'SSL模式',
    'disable': '禁用',
    'require': '必需',
    'verify_ca': '验证CA',
    'verify_full': '完全验证'
  }
};

// Get translation for a key
function __(key) {
  if (translations[currentLanguage] && translations[currentLanguage][key]) {
    return translations[currentLanguage][key];
  }
  // Fallback to English if key not found in current language
  if (translations['en'] && translations['en'][key]) {
    return translations['en'][key];
  }
  // Return the key itself if not found in any language
  return key;
}

// Change language
function changeLanguage(lang) {
  if (availableLanguages[lang]) {
    currentLanguage = lang;
    localStorage.setItem('language', lang);
    // Update UI with new language
    updatePageLanguage();
    return true;
  }
  return false;
}

// Get current language
function getCurrentLanguage() {
  return currentLanguage;
}

// Get available languages
function getAvailableLanguages() {
  return availableLanguages;
}

// Update all text elements on page with translated content
function updatePageLanguage() {
  // Update page title
  const pageTitle = document.querySelector('meta[name="page-title"]');
  if (pageTitle) {
    const titleKey = pageTitle.getAttribute('data-i18n');
    if (titleKey) {
      document.title = __(titleKey) + ' | ' + __('notification_service');
    }
  }
  
  // Update all elements with data-i18n attribute
  document.querySelectorAll('[data-i18n]').forEach(element => {
    const key = element.getAttribute('data-i18n');
    if (key) {
      element.textContent = __(key);
    }
  });
  
  // Update all elements with data-i18n-placeholder attribute
  document.querySelectorAll('[data-i18n-placeholder]').forEach(element => {
    const key = element.getAttribute('data-i18n-placeholder');
    if (key) {
      element.setAttribute('placeholder', __(key));
    }
  });
  
  // Update all elements with data-i18n-title attribute
  document.querySelectorAll('[data-i18n-title]').forEach(element => {
    const key = element.getAttribute('data-i18n-title');
    if (key) {
      element.setAttribute('title', __(key));
    }
  });
  
  // Ensure status badges are translated
  document.querySelectorAll('.badge').forEach(badge => {
    const text = badge.textContent.trim().toLowerCase();
    if (['sent', 'failed', 'pending', 'processing'].includes(text)) {
      badge.textContent = __(text);
    }
  });
  
  // Trigger custom event for components that need to update
  document.dispatchEvent(new CustomEvent('languageChanged', { detail: { language: currentLanguage } }));
}

// Initialize language on page load
document.addEventListener('DOMContentLoaded', function() {
  // Initialize language from local storage or browser preference
  if (!localStorage.getItem('language')) {
    // Try to get browser language
    const browserLang = navigator.language || navigator.userLanguage;
    if (browserLang && browserLang.startsWith('zh')) {
      currentLanguage = 'zh';
    } else {
      currentLanguage = 'en';
    }
    localStorage.setItem('language', currentLanguage);
  }
  
  // Apply initial language
  updatePageLanguage();
}); 