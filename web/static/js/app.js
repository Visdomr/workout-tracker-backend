// Workout Tracker JavaScript

// Global variables
let currentWorkout = null;
let currentExercise = null;

// DOM Content Loaded
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
});

function initializeApp() {
    // Initialize modals
    initializeModals();
    
    // Initialize forms
    initializeForms();
    
    // Initialize workout functionality
    initializeWorkoutFeatures();
    
    // Initialize dark mode
    initializeDarkMode();
    
    console.log('Workout Tracker initialized');
}

// Modal functionality
function initializeModals() {
    // Get all modals
    const modals = document.querySelectorAll('.modal');
    
    modals.forEach(modal => {
        const closeBtn = modal.querySelector('.close');
        
        // Close modal when clicking X
        if (closeBtn) {
            closeBtn.addEventListener('click', function() {
                modal.style.display = 'none';
            });
        }
        
        // Close modal when clicking outside
        window.addEventListener('click', function(event) {
            if (event.target === modal) {
                modal.style.display = 'none';
            }
        });
    });
}

// Form functionality
function initializeForms() {
    // Set today's date as default for date inputs
    const dateInputs = document.querySelectorAll('input[type="date"]');
    const today = new Date().toISOString().split('T')[0];
    
    dateInputs.forEach(input => {
        if (!input.value) {
            input.value = today;
        }
    });
    
    // Don't interfere with any form submissions
    // Let all forms submit normally
}

// Workout features
function initializeWorkoutFeatures() {
    // Only handle specific workout feature buttons, not form submit buttons
    // Don't interfere with form submissions
    console.log('Workout features initialized (not interfering with forms)');
}

// Dark mode functionality
function initializeDarkMode() {
    // Apply dark mode based on user settings from server
    const themeSelect = document.getElementById('theme');
    if (themeSelect) {
        // Apply current theme on page load
        const currentTheme = themeSelect.value;
        applyTheme(currentTheme);
        
        // Listen for theme changes
        themeSelect.addEventListener('change', function() {
            const selectedTheme = this.value;
            applyTheme(selectedTheme);
        });
    }
}

function applyTheme(theme) {
    const body = document.body;
    if (theme === 'dark') {
        body.classList.add('dark-mode');
    } else {
        body.classList.remove('dark-mode');
    }
}

// Modal functions
function openCreateWorkoutModal() {
    const modal = document.getElementById('createWorkoutModal');
    if (modal) {
        modal.style.display = 'block';
    }
}

function closeCreateWorkoutModal() {
    const modal = document.getElementById('createWorkoutModal');
    if (modal) {
        modal.style.display = 'none';
    }
}

function openAddExerciseModal() {
    const modal = document.getElementById('addExerciseModal');
    if (modal) {
        modal.style.display = 'block';
    }
}

function closeAddExerciseModal() {
    const modal = document.getElementById('addExerciseModal');
    if (modal) {
        modal.style.display = 'none';
    }
}

// Form handlers - removed to allow normal form submissions
// Forms will now submit normally without JavaScript interference

// Utility functions
function showMessage(message, type = 'info') {
    // Create a temporary message element
    const messageEl = document.createElement('div');
    messageEl.className = `message message-${type}`;
    messageEl.textContent = message;
    
    // Style the message
    messageEl.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 1rem 1.5rem;
        background: ${type === 'success' ? '#27ae60' : type === 'error' ? '#e74c3c' : '#3498db'};
        color: white;
        border-radius: 4px;
        z-index: 1001;
        animation: slideIn 0.3s ease;
    `;
    
    document.body.appendChild(messageEl);
    
    // Remove after 3 seconds
    setTimeout(() => {
        messageEl.remove();
    }, 3000);
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
}

function formatTime(seconds) {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;
    
    if (hours > 0) {
        return `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
    } else {
        return `${minutes}:${secs.toString().padStart(2, '0')}`;
    }
}

// API helper functions (for future use)
async function apiRequest(url, options = {}) {
    const defaultOptions = {
        headers: {
            'Content-Type': 'application/json',
        },
    };
    
    const response = await fetch(url, { ...defaultOptions, ...options });
    
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    return await response.json();
}

async function createWorkoutAPI(workoutData) {
    try {
        const response = await apiRequest('/api/workouts', {
            method: 'POST',
            body: JSON.stringify(workoutData),
        });
        return response;
    } catch (error) {
        console.error('Error creating workout:', error);
        throw error;
    }
}

async function getWorkoutsAPI() {
    try {
        const response = await apiRequest('/api/workouts');
        return response;
    } catch (error) {
        console.error('Error fetching workouts:', error);
        throw error;
    }
}

// Export functions for use in other scripts
window.WorkoutTracker = {
    openCreateWorkoutModal,
    closeCreateWorkoutModal,
    openAddExerciseModal,
    closeAddExerciseModal,
    showMessage,
    formatDate,
    formatTime,
    apiRequest,
    createWorkoutAPI,
    getWorkoutsAPI
};

// Add CSS for message animations
const style = document.createElement('style');
style.textContent = `
    @keyframes slideIn {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }
    
    .message {
        animation: slideIn 0.3s ease;
    }
`;
document.head.appendChild(style);
