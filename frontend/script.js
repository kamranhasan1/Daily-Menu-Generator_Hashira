document.addEventListener('DOMContentLoaded', function() {
    const dateInput = document.getElementById('date-input');
    const daySelect = document.getElementById('day-select');
    const generateBtn = document.getElementById('generate-btn');
    const loadingElement = document.getElementById('loading');
    const menuContainer = document.getElementById('menu-container');
    const threeDayBtn = document.getElementById('3day-btn');
    const sevenDayBtn = document.getElementById('7day-btn');
    
    let selectedDays = 3;
    
    // Set default date to today
    dateInput.valueAsDate = new Date();
    
    // Duration toggle
    threeDayBtn.addEventListener('click', function() {
        selectedDays = 3;
        threeDayBtn.classList.add('active');
        sevenDayBtn.classList.remove('active');
    });
    
    sevenDayBtn.addEventListener('click', function() {
        selectedDays = 7;
        sevenDayBtn.classList.add('active');
        threeDayBtn.classList.remove('active');
    });
    
    generateBtn.addEventListener('click', function() {
        const selectedDay = daySelect.value;
        
        if (selectedDay) {
            const date = getDateForDay(selectedDay);
            fetchMenu(date);
        } else {
            const selectedDate = dateInput.value;
            if (!selectedDate) {
                alert('Please select either a date or a day');
                return;
            }
            fetchMenu(selectedDate);
        }
    });
    
    function getDateForDay(dayName) {
        const days = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
        const today = new Date();
        const todayDay = today.getDay();
        const targetDay = days.indexOf(dayName);
        
        let diff = targetDay - todayDay;
        if (diff < 0) diff += 7;
        
        const targetDate = new Date(today);
        targetDate.setDate(today.getDate() + diff);
        return targetDate.toISOString().split('T')[0];
    }
    
    function fetchMenu(startDate) {
        loadingElement.classList.remove('hidden');
        menuContainer.innerHTML = '';
        
        fetch(`http://localhost:8080/menu?date=${startDate}&days=${selectedDays}`)
            .then(response => {
                if (!response.ok) throw new Error('Network response was not ok');
                return response.json();
            })
            .then(menuData => {
                loadingElement.classList.add('hidden');
                console.log("API Response:", menuData); // Debug log
                displayMenu(menuData);
            })
            .catch(error => {
                loadingElement.classList.add('hidden');
                menuContainer.innerHTML = `<div class="error">Error: ${error.message}</div>`;
                console.error('Fetch Error:', error);
            });
    }
    
    function displayMenu(menuData) {
        if (!menuData || !Array.isArray(menuData)) {
            menuContainer.innerHTML = '<div class="no-menu">Invalid menu data received</div>';
            return;
        }
        
        menuContainer.innerHTML = '';
        
        menuData.forEach(day => {
            if (!day || !day.meal_options || !Array.isArray(day.meal_options)) {
                console.warn("Invalid day data:", day);
                return;
            }
            
            const dayElement = document.createElement('div');
            dayElement.className = 'day-menu';
            
            const dateHeader = document.createElement('h2');
            dateHeader.textContent = new Date(day.date).toLocaleDateString('en-US', { 
                weekday: 'long', 
                year: 'numeric', 
                month: 'long', 
                day: 'numeric' 
            });
            dayElement.appendChild(dateHeader);
            
            const optionsContainer = document.createElement('div');
            optionsContainer.className = 'meal-options';
            
            day.meal_options.forEach(option => {
                const mealElement = document.createElement('div');
                mealElement.className = 'meal-option';
                
                mealElement.innerHTML = `
                    <h3>${option.main?.name || 'Not specified'}</h3>
                    <p><strong>Side:</strong> ${option.side?.name || 'Not specified'}</p>
                    <p><strong>Drink:</strong> ${option.drink?.name || 'Not specified'}</p>
                    <div class="nutrition-info">
                        <span>🔥 ${option.total_calories || 0} calories</span>
                        <span>⭐ ${option.combined_popularity?.toFixed(1) || '0.0'} popularity</span>
                    </div>
                    <div class="taste-profile">
                        Flavor: ${option.main?.taste_profile || 'unknown'}/
                        ${option.side?.taste_profile || 'unknown'}/
                        ${option.drink?.taste_profile || 'unknown'}
                    </div>
                `;
                
                optionsContainer.appendChild(mealElement);
            });
            
            dayElement.appendChild(optionsContainer);
            menuContainer.appendChild(dayElement);
        });
    }
});