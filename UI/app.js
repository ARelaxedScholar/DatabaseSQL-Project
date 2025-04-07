// app.js - Adjusted to use DTOs and match DB schema

document.addEventListener('DOMContentLoaded', () => {
    // --- Constants and State ---
    const API_BASE_URL = 'https://sunflower-booking-backend-966219880837.us-central1.run.app';

    const views = document.querySelectorAll('.view');
    const navLinks = document.querySelectorAll('.nav-links a[data-view]');
    const feedbackElements = document.querySelectorAll('.feedback');

    // --- Helper Functions ---
    // decodeJwt, displayFeedback, clearAllFeedback, apiRequest (modified as needed)

    /**
     * Decodes a JWT token to extract payload data (basic implementation).
     * WARNING: This does NOT verify the signature. Verification MUST happen server-side.
     * This is only for conveniently reading claims client-side after server validation.
     * @param {string} token - The JWT token.
     * @returns {object|null} - The decoded payload or null if decoding fails.
     */
    function decodeJwt(token) {
        try {
            const base64Url = token.split('.')[1];
            const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
            const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
                return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
            }).join(''));
            return JSON.parse(jsonPayload);
        } catch (error) {
            console.error("Failed to decode JWT:", error);
            return null;
        }
    }

    /**
     * Displays feedback messages in designated areas.
     * @param {string} elementId - The ID of the feedback <p> element.
     * @param {string} message - The message to display.
     * @param {boolean} isError - True for error styling, false for success.
     */
    function displayFeedback(elementId, message, isError = false) {
        const el = document.getElementById(elementId);
        if (!el) {
            console.error(`Feedback element not found: ${elementId}`);
            alert(message);
            return;
        }
        el.classList.remove('success', 'error');
        el.textContent = message;
        el.classList.add(isError ? 'error' : 'success');
        el.style.display = 'block';
    }

    /** Clears all feedback messages */
    function clearAllFeedback() {
        feedbackElements.forEach(el => {
            el.textContent = '';
            el.style.display = 'none';
            el.classList.remove('success', 'error');
        });
    }

    /**
     * Centralized function for making API requests using fetch.
     * @param {string} endpoint - The API endpoint.
     * @param {string} method - HTTP method (GET, POST, PUT, DELETE, PATCH).
     * @param {object|null} body - Data to send in the request body.
     * @param {boolean} requiresAuth - Whether to include the JWT Authorization header.
     * @returns {Promise<any>} - Resolves with the JSON response data.
     */
    async function apiRequest(endpoint, method = 'GET', body = null, requiresAuth = false) {
        const url = `${API_BASE_URL}${endpoint}`;
        const options = {
            method: method,
            headers: {
                'Accept': 'application/json'
            }
        };

        if (requiresAuth) {
            const token = localStorage.getItem('jwt');
            if (!token) {
                console.error('Auth token not found.');
                showView('login-view');
                throw new Error('Veuillez vous reconnecter.');
            }
            options.headers['Authorization'] = `Bearer ${token}`;
        }

        if (body && (method === 'POST' || method === 'PUT' || method === 'PATCH')) {
            options.headers['Content-Type'] = 'application/json';
            options.body = JSON.stringify(body);
        }

        try {
            const response = await fetch(url, options);
            let responseData = null;
            const contentType = response.headers.get('content-type');

            if (response.status === 204) {
                return null;
            }

            if (contentType && contentType.includes('application/json')) {
                try {
                    responseData = await response.json();
                } catch (jsonError) {
                    if (!response.ok) {
                        console.warn("Could not parse JSON error response body", jsonError);
                    }
                }
            }

            if (!response.ok) {
                const errorMessage = responseData?.message || responseData?.error || `HTTP error ${response.status}`;
                console.error(`API Error (${response.status}) for ${method} ${endpoint}: ${errorMessage}`, responseData);
                throw new Error(errorMessage);
            }

            return responseData;

        } catch (error) {
            console.error('API request failed:', error);
            if (error instanceof TypeError) {
                throw new Error(`Erreur réseau ou serveur indisponible.`);
            }
            throw error;
        }
    }

    /** Checks if the logged-in user has the admin flag set */
    function isAdminUser() {
        return localStorage.getItem('isAdmin') === 'true';
    }

    /** Updates navigation links based on login status and role */
    function updateNav() {
        const token = localStorage.getItem('jwt');
        const role = localStorage.getItem('role');
        const isAdmin = isAdminUser();

        const loginNav = document.getElementById('login-nav');
        const registerNav = document.getElementById('register-nav');
        const clientProfileNav = document.getElementById('client-profile-nav');
        const clientReservationsNav = document.getElementById('client-reservations-nav');
        const employeeDashboardNav = document.getElementById('employee-dashboard-nav');
        const adminDashboardNav = document.getElementById('admin-dashboard-nav');
        const viewsNav = document.getElementById('views-nav');
        const logoutNav = document.getElementById('logout-nav');

        [clientProfileNav, clientReservationsNav, employeeDashboardNav, adminDashboardNav, viewsNav, logoutNav].forEach(el => el.style.display = 'none');
        [loginNav, registerNav].forEach(el => el.style.display = 'list-item');

        if (token) {
            loginNav.style.display = 'none';
            registerNav.style.display = 'none';
            logoutNav.style.display = 'list-item';

            if (role === 'client') {
                clientProfileNav.style.display = 'list-item';
                clientReservationsNav.style.display = 'list-item';
            } else if (role === 'employee') {
                employeeDashboardNav.style.display = 'list-item';
                viewsNav.style.display = 'list-item';
                if (isAdmin) {
                    adminDashboardNav.style.display = 'list-item';
                }
            }
        }
    }

    /** Shows a specific view section and hides others. Also triggers data loading. */
    function showView(viewId) {
        clearAllFeedback();
        views.forEach(view => view.classList.remove('active'));
        navLinks.forEach(link => {
            link.classList.remove('active-nav');
            if (link.getAttribute('href') === `#${viewId}` || link.dataset.view === viewId) {
                link.classList.add('active-nav');
            }
        });

        const viewToShow = document.getElementById(viewId);
        if (viewToShow) {
            viewToShow.classList.add('active');
            switch (viewId) {
                case 'search-view': loadSearchDropdowns(); break;
                case 'client-profile-view': loadClientProfile(); break;
                case 'client-reservations-view': loadClientReservations(); break;
                case 'admin-dashboard-view': if(isAdminUser()) loadAdminData(); else console.warn("Attempt to load admin data without admin rights"); break;
                case 'required-views-view': if(localStorage.getItem('role') === 'employee') loadRequiredViewsData(); else console.warn("Attempt to load views data without employee rights"); break;
                case 'employee-dashboard-view': /* Load data if needed */ break;
            }
        } else {
            console.warn(`View with ID "${viewId}" not found.`);
            document.getElementById('search-view')?.classList.add('active');
            document.querySelector('.nav-links a[data-view="search-view"]')?.classList.add('active-nav');
        }
        updateNav();
    }

    // --- Authentication Logic ---
    function handleMagicLink() {
        const params = new URLSearchParams(window.location.search);
        const token = params.get('token');
        const role = params.get('role');
        const isAdminParam = params.get('admin');

        if (token && role) {
            localStorage.setItem('jwt', token);
            localStorage.setItem('role', role);

            if (isAdminParam === 'true' && role === 'employee') {
                localStorage.setItem('isAdmin', 'true');
                console.log("Admin flag set.");
            } else {
                localStorage.removeItem('isAdmin');
            }

            history.replaceState(null, '', window.location.pathname);
            console.log(`Logged in as ${role}${isAdminUser() ? ' (Admin)' : ''}`);

            if (role === 'client') {
                showView('search-view');
            } else if (role === 'employee') {
                if (isAdminUser()) {
                    showView('admin-dashboard-view');
                } else {
                    showView('employee-dashboard-view');
                }
            }
            updateNav();
            return true;
        }
        return false;
    }

    function logout() {
        localStorage.removeItem('jwt');
        localStorage.removeItem('role');
        localStorage.removeItem('isAdmin');
        localStorage.removeItem('clientId'); // clear stored client ID
        console.log("Logged out");
        updateNav();
        showView('login-view');
    }

    // --- Event Listeners ---
    navLinks.forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const viewId = link.dataset.view || link.getAttribute('href').substring(1);
            const role = localStorage.getItem('role');
            const token = localStorage.getItem('jwt');
            const isAdmin = isAdminUser();

            if (!token && (viewId.startsWith('client-') || viewId.startsWith('employee-') || viewId.startsWith('admin-') || viewId.startsWith('required-views'))) {
                showView('login-view');
                displayFeedback('login-feedback', 'Veuillez vous connecter pour accéder à cette page.', true);
                return;
            }

            let authorized = true;
            if (viewId.startsWith('client-') && role !== 'client') authorized = false;
            if (viewId.startsWith('employee-') && role !== 'employee') authorized = false;
            if (viewId.startsWith('admin-') && !isAdmin) authorized = false;
            if (viewId.startsWith('required-views-') && role !== 'employee') authorized = false;

            if (authorized) {
                showView(viewId);
            } else {
                const currentActiveView = document.querySelector('.view.active');
                const feedbackElementId = currentActiveView ? `${currentActiveView.id}-feedback` : 'login-feedback';
                displayFeedback(feedbackElementId, 'Accès non autorisé pour votre rôle.', true);
                console.warn(`Unauthorized attempt to access ${viewId} by role: ${role}, isAdmin: ${isAdmin}`);
            }
        });
    });

    document.getElementById('logout-button')?.addEventListener('click', (e) => {
        e.preventDefault();
        logout();
    });

    // --- Login Form ---
    document.getElementById('login-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        clearAllFeedback();
        const email = e.target.email.value.trim();
        const role = e.target.role.value;
        const feedbackId = 'login-feedback';
        if (!email) return displayFeedback(feedbackId, 'Adresse e-mail requise.', true);
        const endpoint = role === 'client' ? '/clients/login' : '/employees/login';
        try {
            await apiRequest(endpoint, 'POST', { email });
            displayFeedback(feedbackId, 'Lien de connexion envoyé ! Veuillez vérifier votre boîte e-mail.', false);
            e.target.reset();
        } catch (error) {
            displayFeedback(feedbackId, `Erreur de connexion : ${error.message}`, true);
        }
    });

    // --- Registration Form ---
    document.getElementById('register-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        clearAllFeedback();
        const feedbackId = 'register-feedback';
        const formData = {
            sin: e.target.sin.value.trim(),
            firstName: e.target.firstName.value.trim(),
            lastName: e.target.lastName.value.trim(),
            address: e.target.address.value.trim(),
            phone: e.target.phone.value.trim(),
            email: e.target.email.value.trim(),
            joinDate: new Date().toISOString()
        };
        if (!formData.sin || !formData.firstName || !formData.lastName || !formData.email) {
            return displayFeedback(feedbackId, 'Veuillez remplir tous les champs requis.', true);
        }
        if (formData.sin.length !== 9 || !/^\d{9}$/.test(formData.sin)) {
            return displayFeedback(feedbackId, 'Le SIN doit contenir exactement 9 chiffres.', true);
        }
        try {
            const result = await apiRequest('/clients/register', 'POST', formData);
            displayFeedback(feedbackId, `Inscription réussie ! Client ID: ${result?.clientId ?? 'N/A'}. Vous pouvez maintenant vous connecter.`, false);
            e.target.reset();
        } catch (error) {
            displayFeedback(feedbackId, `Erreur d'inscription : ${error.message}`, true);
        }
    });

    // --- Search/Booking Logic ---
    // Global mapping object so we can convert from room type ID to name later
let roomTypeMapping = {};

/**
 * Loads room types from the backend and populates the "search-room-type" select.
 * Each option's value is the room type's numeric ID, and its text is "ID: Name".
 */
async function loadSearchRoomTypes() {
    try {
        const roomTypes = await apiRequest('/roomtypes', 'GET', null, false);
        if (roomTypes && Array.isArray(roomTypes)) {
            // Sort room types by ID (ascending)
            roomTypes.sort((a, b) => a.id - b.id);
            // Build mapping and populate select
            roomTypeMapping = {};
            const select = document.getElementById('search-room-type');
            select.innerHTML = `<option value="">Tous</option>`;
            roomTypes.forEach(rt => {
                roomTypeMapping[rt.id] = rt.name;
                const option = document.createElement('option');
                option.value = rt.id; // Numeric ID as value
                option.textContent = `${rt.id}: ${rt.name}`;
                select.appendChild(option);
            });
        }
    } catch (error) {
        console.error("Failed to load room types for search:", error);
    }
}


    const searchForm = document.getElementById('search-form');
    const roomResultsContainer = document.getElementById('room-results-container');
    let lastSearchParams = {};

    async function loadSearchDropdowns() {
        // Load Hotel Chains
        const chainSelect = document.getElementById('search-hotel-chain');
        if (chainSelect && chainSelect.options.length <= 1) {
            try {
                const chains = await apiRequest('/admin/hotelchains', 'GET', null, false);
                if (chains) {
                    populateSelect('search-hotel-chain', chains, 'id', 'name', 'Toutes');
                }
            } catch (error) {
                console.error("Failed to load hotel chains for search:", error);
            }
            await loadSearchRoomTypes();
        }
        // Load Room Types
        const roomTypeSelect = document.getElementById('search-room-type');
        if (roomTypeSelect && roomTypeSelect.options.length <= 1) {
            try {
                const types = await apiRequest('/roomtypes', 'GET', null, false);
                if (types) {
                    // Use room type name as value (per DTO: RoomType is a string)
                    populateSelect('search-room-type', types, 'name', 'name', 'Tous');
                }
            } catch (error) {
                console.error("Failed to load room types for search:", error);
            }
        }
    }

    // Helper to populate a select element (assumes data is an array of objects)
    function populateSelect(selectElementId, data, valueField = 'id', textField = 'name', prompt = 'Sélectionner...') {
        const select = document.getElementById(selectElementId);
        if (!select || !Array.isArray(data)) return;
        select.innerHTML = `<option value="">${prompt}</option>`;
        data.forEach(item => {
            if(item && item[valueField] !== undefined && item[textField] !== undefined) {
                const option = document.createElement('option');
                option.value = item[valueField];
                option.textContent = item[textField];
                select.appendChild(option);
            }
        });
    }

    function formatDateForBackend(dateStr) {
        const dateObj = new Date(dateStr);
        const month = (dateObj.getMonth() + 1).toString().padStart(2, '0');
        const day = dateObj.getDate().toString().padStart(2, '0');
        const year = dateObj.getFullYear().toString(); // Full year
        return `${month}-${day}-${year}`;
      }
      

      searchForm?.addEventListener('submit', async (e) => {
        e.preventDefault();
        clearAllFeedback();
        const feedbackId = 'search-feedback';
        roomResultsContainer.innerHTML = '<p>Recherche en cours...</p>';
    
        const formData = new FormData(e.target);
        const rawStartDate = formData.get('startDate');
        const rawEndDate = formData.get('endDate');
    
        if (!rawStartDate || !rawEndDate) {
            roomResultsContainer.innerHTML = '';
            return displayFeedback(feedbackId, 'Les dates d\'arrivée et de départ sont requises.', true);
        }
        if (new Date(rawEndDate) <= new Date(rawStartDate)) {
            roomResultsContainer.innerHTML = '';
            return displayFeedback(feedbackId, 'La date de départ doit être après la date d\'arrivée.', true);
        }
    
        // Convert dates using your helper (assume MM-DD-YYYY as required)
        const startDate = formatDateForBackend(rawStartDate);
        const endDate = formatDateForBackend(rawEndDate);
    
        // Convert optional numeric fields
        const capacity = formData.get('capacity') ? parseInt(formData.get('capacity')) : undefined;
        const priceMin = formData.get('priceMin') ? parseFloat(formData.get('priceMin')) : undefined;
        const priceMax = formData.get('priceMax') ? parseFloat(formData.get('priceMax')) : undefined;
        const hotelChainId = formData.get('hotelChainId') ? parseInt(formData.get('hotelChainId')) : undefined;
    
        // For room type, the select returns a numeric ID (as a string). If provided, map it to the corresponding name.
        const roomTypeId = formData.get('roomType');
        let roomType;
        if (roomTypeId) {
            roomType = roomTypeMapping[roomTypeId];
        }
    
        // Build the payload object following RoomSearchInput DTO
        const searchPayload = {
            startDate: startDate,
            endDate: endDate,
            capacity: capacity,
            priceMin: priceMin,
            priceMax: priceMax,
            hotelChainId: hotelChainId,
            roomType: roomType
        };
    
        // Create URL query parameters, omitting undefined values
        const params = new URLSearchParams();
        for (const key in searchPayload) {
            if (searchPayload[key] !== undefined && searchPayload[key] !== '') {
                params.append(key, searchPayload[key]);
            }
        }
    
        try {
            const searchResult = await apiRequest(`/search/rooms?${params.toString()}`, 'GET', null, false);
            const rooms = searchResult?.rooms;
            if (!rooms || rooms.length === 0) {
                roomResultsContainer.innerHTML = '<p>Aucune chambre disponible pour les critères sélectionnés.</p>';
            } else {
                roomResultsContainer.innerHTML = rooms.map(room => `
                    <div class="room-result" data-room-id="${room.roomId}" data-hotel-id="${room.hotelId}" data-price="${room.price}">
                        <div class="details">
                            <h4>Chambre ${room.number || 'N/A'} - Étage ${room.floor || 'N/A'} (Hôtel ID: ${room.hotelId})</h4>
                            <p><strong>Type:</strong> ${room.roomType || 'N/A'}</p>
                            <p><strong>Capacité:</strong> ${room.capacity ?? '?'} | <strong>Superficie:</strong> ${room.surfaceArea ?? '?'} m²</p>
                            <p><strong>Prix:</strong> ${(room.price ?? 0).toFixed(2)} € / nuit</p>
                            <p><strong>Aménités:</strong> ${Array.isArray(room.amenities) ? room.amenities.join(', ') : 'N/A'}</p>
                            <p><strong>Vues:</strong> ${Array.isArray(room.viewTypes) ? room.viewTypes.join(', ') : 'N/A'}</p>
                            <p><strong>Extensible:</strong> ${room.isExtensible ? 'Oui' : 'Non'}</p>
                        </div>
                        <button class="reserve-button" type="button">Réserver</button>
                    </div>
                `).join('');
            }
        } catch (error) {
            displayFeedback(feedbackId, `Erreur de recherche : ${error.message}`, true);
            roomResultsContainer.innerHTML = '<p>Erreur lors de la recherche.</p>';
        }
    });
    

    // Reservation event: compute totalPrice and include DTO fields
    roomResultsContainer?.addEventListener('click', async (e) => {
        if (e.target.classList.contains('reserve-button')) {
            e.preventDefault();
            clearAllFeedback();
            const roomResultDiv = e.target.closest('.room-result');
            const roomId = roomResultDiv?.dataset.roomId;
            const hotelId = roomResultDiv?.dataset.hotelId;
            const feedbackId = 'search-feedback';

            if (!localStorage.getItem('jwt') || localStorage.getItem('role') !== 'client') {
                displayFeedback(feedbackId, 'Veuillez vous connecter en tant que client pour réserver.', true);
                showView('login-view');
                return;
            }
            if (!roomId || !hotelId || !lastSearchParams.startDate || !lastSearchParams.endDate) {
                displayFeedback(feedbackId, 'Information de réservation manquante. Veuillez relancer la recherche.', true);
                return;
            }
            if (!confirm(`Confirmer la réservation pour la chambre ID ${roomId} du ${lastSearchParams.startDate} au ${lastSearchParams.endDate}?`)) return;

            // Compute total price: number of nights * room price
            const start = new Date(lastSearchParams.startDate);
            const end = new Date(lastSearchParams.endDate);
            const nights = Math.ceil((end - start) / (1000 * 60 * 60 * 24));
            const price = parseFloat(roomResultDiv.dataset.price);
            const totalPrice = price * nights;

            // Retrieve clientId (stored from profile load)
            const clientId = parseInt(localStorage.getItem('clientId'), 10);

            const reservationData = {
                clientId: clientId,
                hotelID: parseInt(hotelId, 10), // Note: key "hotelID" as per DTO
                roomId: parseInt(roomId, 10),
                startDate: lastSearchParams.startDate,
                endDate: lastSearchParams.endDate,
                reservationDate: new Date().toISOString(),
                totalPrice: totalPrice,
                status: 1 // Assuming 1 = "Confirmed"
            };

            try {
                const result = await apiRequest('/clients/reservations', 'POST', reservationData, true);
                displayFeedback(feedbackId, `Réservation réussie ! ID : ${result.reservationId}.`, false);
            } catch (error) {
                displayFeedback(feedbackId, `Erreur de réservation : ${error.message}`, true);
            }
        }
    });

    // --- Client Profile Logic ---
    const profileForm = document.getElementById('client-profile-form');

    async function loadClientProfile() {
        clearAllFeedback();
        const feedbackId = 'profile-feedback';
        try {
            const profileData = await apiRequest('/clients/profile', 'GET', null, true);
            // Store clientId locally for use in updates and reservations
            localStorage.setItem('clientId', profileData.clientId);
            document.getElementById('profile-sin').textContent = profileData.sin || 'N/A';
            profileForm.elements['firstName'].value = profileData.firstName || '';
            profileForm.elements['lastName'].value = profileData.lastName || '';
            profileForm.elements['address'].value = profileData.address || '';
            profileForm.elements['phone'].value = profileData.phone || '';
            profileForm.elements['email'].value = profileData.email || '';
            document.getElementById('profile-joindate').textContent = profileData.joinDate ? new Date(profileData.joinDate).toLocaleDateString() : 'N/A';
        } catch (error) {
            displayFeedback(feedbackId, `Erreur chargement profil : ${error.message}`, true);
            logout();
        }
    }

    profileForm?.addEventListener('submit', async (e) => {
        e.preventDefault();
        clearAllFeedback();
        const feedbackId = 'profile-feedback';
        const clientId = parseInt(localStorage.getItem('clientId'), 10);
        const formData = {
            clientId: clientId,
            firstName: e.target.elements['firstName'].value.trim(),
            lastName: e.target.elements['lastName'].value.trim(),
            address: e.target.elements['address'].value.trim(),
            phone: e.target.elements['phone'].value.trim(),
            email: e.target.elements['email'].value.trim(),
        };
        if (!formData.firstName || !formData.lastName || !formData.email) {
            displayFeedback(feedbackId, 'Prénom, Nom et Email sont requis.', true);
            return;
        }
        try {
            await apiRequest('/clients/profile', 'PUT', formData, true);
            displayFeedback(feedbackId, 'Profil mis à jour.', false);
        } catch (error) {
            displayFeedback(feedbackId, `Erreur mise à jour : ${error.message}`, true);
        }
    });

    // --- Client Reservations Logic ---
    const reservationsListDiv = document.getElementById('client-reservations-list');

    function mapReservationStatus(statusInt) {
        switch(statusInt) {
            case 1: return 'Confirmée';
            case 2: return 'En attente';
            case 3: return 'Annulée';
            case 4: return 'Terminée';
            default: return `Inconnu (${statusInt})`;
        }
    }

    async function loadClientReservations() {
        clearAllFeedback();
        reservationsListDiv.innerHTML = '<p>Chargement...</p>';
        const feedbackId = 'reservations-feedback';
        try {
            const reservations = await apiRequest('/clients/reservations', 'GET', null, true);
            if (!reservations || reservations.length === 0) {
                reservationsListDiv.innerHTML = '<p>Aucune réservation.</p>';
            } else {
                reservationsListDiv.innerHTML = reservations.map(res => `
                    <div class="reservation-item" data-reservation-id="${res.reservationId}">
                        <p><strong>Réservation ID:</strong> ${res.reservationId}</p>
                        <p><strong>Hôtel ID:</strong> ${res.hotelID} | <strong>Chambre ID:</strong> ${res.roomId}</p>
                        <p><strong>Dates:</strong> ${new Date(res.startDate).toLocaleDateString()} - ${new Date(res.endDate).toLocaleDateString()}</p>
                        <p><strong>Prix Total:</strong> ${res.totalPrice?.toFixed(2) ?? 'N/A'} €</p>
                        <p><strong>Statut:</strong> ${mapReservationStatus(res.status)}</p>
                        <p><strong>Date Réservation:</strong> ${new Date(res.reservationDate).toLocaleString()}</p>
                        ${(res.status === 1 || res.status === 2) ? 
                           '<button class="cancel-reservation-button delete-button" type="button">Annuler</button>' : '' }
                    </div>
                `).join('');
            }
        } catch (error) {
            displayFeedback(feedbackId, `Erreur chargement réservations : ${error.message}`, true);
            reservationsListDiv.innerHTML = '<p>Impossible de charger.</p>';
        }
    }

    reservationsListDiv?.addEventListener('click', async (e) => {
        if (e.target.classList.contains('cancel-reservation-button')) {
            e.preventDefault(); 
            clearAllFeedback();
            const reservationId = e.target.closest('.reservation-item')?.dataset.reservationId;
            const feedbackId = 'reservations-feedback';
            if (!reservationId) return;
            if (!confirm(`Annuler la réservation ID ${reservationId} ?`)) return;
            try {
                await apiRequest(`/clients/reservations/${reservationId}`, 'DELETE', null, true);
                displayFeedback(feedbackId, `Réservation ID ${reservationId} annulée.`, false);
                loadClientReservations();
            } catch (error) {
                displayFeedback(feedbackId, `Erreur annulation : ${error.message}`, true);
            }
        }
    });

    // --- Employee Dashboard Logic ---
    // Check-in Form
    document.getElementById('employee-checkin-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        clearAllFeedback();
        const feedbackId = 'checkin-feedback';
        const reservationIdStr = e.target.elements['reservationId'].value.trim();
        if (!reservationIdStr || isNaN(parseInt(reservationIdStr))) {
            return displayFeedback(feedbackId, 'ID réservation valide requis.', true);
        }
        const reservationId = parseInt(reservationIdStr);
        const employeeData = decodeJwt(localStorage.getItem('jwt'));
        const employeeId = employeeData ? employeeData.employeeId : null;
        if (!employeeId) return displayFeedback(feedbackId, 'Employé introuvable.', true);
        try {
            const payload = {
                reservationId: reservationId,
                employeeId: employeeId,
                checkInTime: new Date().toISOString()
            };
            const result = await apiRequest('/employees/checkin', 'POST', payload, true);
            displayFeedback(feedbackId, `Check-in réussi! Séjour ID: ${result.stayId}.`, false);
            e.target.reset();
        } catch (error) {
            displayFeedback(feedbackId, `Erreur check-in : ${error.message}`, true);
        }
    });

    // Create Stay (Walk-in) Form
    document.getElementById('employee-stay-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        clearAllFeedback();
        const feedbackId = 'stay-feedback';
        const employeeData = decodeJwt(localStorage.getItem('jwt'));
        const employeeId = employeeData ? employeeData.employeeId : null;
        if (!employeeId) return displayFeedback(feedbackId, 'Employé introuvable.', true);
        const reservationIDField = e.target.elements['reservationID'];
        const reservationID = reservationIDField && reservationIDField.value.trim() ? parseInt(reservationIDField.value.trim()) : undefined;
        const formData = {
            clientId: parseInt(e.target.elements['clientId'].value.trim()),
            roomId: parseInt(e.target.elements['roomId'].value.trim()),
            // If a reservationID is provided, include it; otherwise, omit for walk-in
            reservationID: reservationID,
            checkInEmployeeId: employeeId,
            checkInTime: e.target.elements['arrivalDate'].value,
            comments: e.target.elements['comments'].value.trim()
        };
        if (isNaN(formData.clientId) || isNaN(formData.roomId) || !formData.checkInTime) {
            return displayFeedback(feedbackId, 'Client ID, Chambre ID et Date d\'arrivée requis.', true);
        }
        try {
            const result = await apiRequest('/employees/stay', 'POST', formData, true);
            displayFeedback(feedbackId, `Location créée! ID: ${result.stayId}.`, false);
            e.target.reset();
        } catch (error) {
            displayFeedback(feedbackId, `Erreur création séjour : ${error.message}`, true);
        }
    });

    // Checkout Form
    document.getElementById('employee-checkout-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        clearAllFeedback();
        const feedbackId = 'checkout-feedback';
        const stayIdStr = e.target.elements['stayId'].value.trim();
        const paymentMethod = e.target.elements['paymentMethod'].value.trim();
        if (!stayIdStr || isNaN(parseInt(stayIdStr)) || !paymentMethod) {
            return displayFeedback(feedbackId, 'ID séjour et méthode paiement requis.', true);
        }
        const employeeData = decodeJwt(localStorage.getItem('jwt'));
        const employeeId = employeeData ? employeeData.employeeId : null;
        if (!employeeId) return displayFeedback(feedbackId, 'Employé introuvable.', true);
        const finalPrice = parseFloat(e.target.elements['finalPrice']?.value) || 0.0;
        const payload = {
            stayID: parseInt(stayIdStr, 10),
            empoyeeID: employeeId,  // Note: using key "empoyeeID" per DTO
            checkOutTime: new Date().toISOString(),
            finalPrice: finalPrice,
            paymentMethod: paymentMethod
        };
        try {
            const result = await apiRequest('/employees/checkout', 'POST', payload, true);
            displayFeedback(feedbackId, `Départ finalisé pour séjour ${payload.stayID}. ${result?.message || ''}`, false);
            e.target.reset();
        } catch (error) {
            displayFeedback(feedbackId, `Erreur départ : ${error.message}`, true);
        }
    });

    // --- Admin Dashboard Logic ---
    let currentEditData = { chain: null, hotel: null, room: null, client: null, employee: null };

    const adminChainForm = document.getElementById('admin-chain-form');
    const adminHotelForm = document.getElementById('admin-hotel-form');
    const adminRoomForm = document.getElementById('admin-room-form');
    const adminClientForm = document.getElementById('admin-client-form');
    const adminEmployeeForm = document.getElementById('admin-employee-form');

    const adminChainsListDiv = document.getElementById('admin-chains-list');
    const adminHotelsListDiv = document.getElementById('admin-hotels-list');
    const adminRoomsListDiv = document.getElementById('admin-rooms-list');
    const adminClientsListDiv = document.getElementById('admin-clients-list');
    const adminEmployeesListDiv = document.getElementById('admin-employees-list');

    // --- populateSelect remains the same (see above) ---

    // --- Admin: Hotel Chain Management ---
    async function loadAdminHotelChains() {
        const feedbackId = 'admin-chain-feedback';
        try {
            const chains = await apiRequest('/admin/hotelchains', 'GET', null, true);
            adminChainsListDiv.innerHTML = createHtmlTable(
                chains || [],
                ['id', 'name', 'numberOfHotels', 'centralAddress', 'email', 'telephone'],
                'chain'
            );
            populateSelect('admin-hotel-filter-chain', chains || []);
            populateSelect('hotel-chain-select', chains || []);
        } catch (error) {
            displayFeedback(feedbackId, `Erreur chargement chaînes: ${error.message}`, true);
            adminChainsListDiv.innerHTML = '<p>Impossible de charger.</p>';
        }
    }
    document.getElementById('add-chain-button')?.addEventListener('click', () => {
        currentEditData.chain = null;
        adminChainForm.reset();
        adminChainForm.style.display = 'block';
        document.getElementById('chain-edit-id').value = '';
    });
    document.getElementById('cancel-chain-button')?.addEventListener('click', () => {
        adminChainForm.style.display = 'none';
        adminChainForm.reset();
        clearAllFeedback();
    });
    adminChainForm?.addEventListener('submit', async (e) => {
        e.preventDefault();
        const feedbackId = 'admin-chain-feedback';
        clearAllFeedback();
        const formData = {
            id: currentEditData.chain ? parseInt(currentEditData.chain.id) : undefined,
            numberOfHotels: 0, // Backend will update this value
            name: e.target.elements['name'].value.trim(),
            centralAddress: e.target.elements['centralAddress'].value.trim(),
            email: e.target.elements['email'].value.trim(),
            telephone: e.target.elements['telephone'].value.trim()
        };
        if (!formData.name || !formData.centralAddress || !formData.email || !formData.telephone) {
            return displayFeedback(feedbackId, 'Champs requis.', true);
        }
        try {
            let result;
            const url = currentEditData.chain ? `/admin/hotelchains/${currentEditData.chain.id}` : '/admin/hotelchains';
            const method = currentEditData.chain ? 'PUT' : 'POST';
            result = await apiRequest(url, method, formData, true);
            displayFeedback(feedbackId, currentEditData.chain ? 'Chaîne màj.' : `Chaîne ajoutée (ID: ${result?.chainId ?? 'N/A'}).`, false);
            adminChainForm.reset();
            adminChainForm.style.display = 'none';
            loadAdminHotelChains();
        } catch (error) {
            displayFeedback(feedbackId, `Erreur sauvegarde: ${error.message}`, true);
        }
    });
    adminChainsListDiv?.addEventListener('click', (e) => {
        if (e.target.classList.contains('edit-button')) {
            const id = e.target.dataset.id;
            const row = e.target.closest('tr');
            currentEditData.chain = {
                id: id,
                name: row.cells[1].textContent,
                numberOfHotels: row.cells[2].textContent,
                centralAddress: row.cells[3].textContent,
                email: row.cells[4].textContent,
                telephone: row.cells[5].textContent,
            };
            adminChainForm.elements['chain-edit-id'].value = id;
            adminChainForm.elements['name'].value = currentEditData.chain.name;
            adminChainForm.elements['centralAddress'].value = currentEditData.chain.centralAddress;
            adminChainForm.elements['email'].value = currentEditData.chain.email;
            adminChainForm.elements['telephone'].value = currentEditData.chain.telephone;
            adminChainForm.style.display = 'block';
            window.scrollTo({ top: adminChainForm.offsetTop - 80, behavior: 'smooth' });
        } else if (e.target.classList.contains('delete-button')) {
            const id = e.target.dataset.id;
            if (confirm(`Supprimer chaîne ID ${id} ? (Supprime aussi hôtels/chambres)`)) {
                deleteAdminItem(`/admin/hotelchains/${id}`, 'admin-chain-feedback', loadAdminHotelChains);
            }
        }
    });

    // --- Admin: Hotel Management ---
    const hotelFilterChain = document.getElementById('admin-hotel-filter-chain');
    async function loadAdminHotels(chainId = null) {
        const feedbackId = 'admin-hotel-feedback';
        let endpoint = '/admin/hotels';
        if (chainId) {
            endpoint += `?chainId=${chainId}`;
            console.warn("Assumed hotel filtering endpoint:", endpoint);
        }
        try {
            const hotels = await apiRequest(endpoint, 'GET', null, true);
            adminHotelsListDiv.innerHTML = createHtmlTable(
                hotels || [],
                ['id', 'chainId', 'rating', 'numberOfRooms', 'name', 'address', 'city', 'email', 'telephone'],
                'hotel'
            );
            populateSelect('admin-room-filter-hotel', hotels || []);
            populateSelect('room-hotel-select', hotels || []);
            populateSelect('employee-acc-hotel', hotels || []);
            populateSelect('view-hotel-capacity-select', hotels || [], 'id', 'name', 'Choisir Hôtel...');
        } catch (error) {
            displayFeedback(feedbackId, `Erreur chargement hôtels: ${error.message}`, true);
            adminHotelsListDiv.innerHTML = '<p>Impossible de charger.</p>';
        }
    }
    hotelFilterChain?.addEventListener('change', (e) => {
        loadAdminHotels(e.target.value || null);
    });
    document.getElementById('add-hotel-button')?.addEventListener('click', () => {
        currentEditData.hotel = null;
        adminHotelForm.reset();
        adminHotelForm.style.display = 'block';
        document.getElementById('hotel-edit-id').value = '';
        if(document.getElementById('hotel-chain-select').options.length <= 1) {
            loadAdminHotelChains().catch(err => console.error("Failed to populate chain select", err));
        }
    });
    document.getElementById('cancel-hotel-button')?.addEventListener('click', () => {
        adminHotelForm.style.display = 'none';
        adminHotelForm.reset();
        clearAllFeedback();
    });
    adminHotelForm?.addEventListener('submit', async (e) => {
        e.preventDefault();
        const feedbackId = 'admin-hotel-feedback';
        clearAllFeedback();
        const formData = {
            id: currentEditData.hotel ? parseInt(currentEditData.hotel.id) : undefined,
            chainId: parseInt(e.target.elements['chainId'].value),
            rating: parseInt(e.target.elements['rating'].value),
            name: e.target.elements['name'].value.trim(),
            address: e.target.elements['address'].value.trim(),
            city: e.target.elements['city'].value.trim(),
            email: e.target.elements['email'].value.trim(),
            phone: e.target.elements['telephone'].value.trim() // Changed key from telephone to phone per DTO
        };
        if (isNaN(formData.chainId) || !formData.name || !formData.address || !formData.city || !formData.email || !formData.phone || isNaN(formData.rating)) {
            return displayFeedback(feedbackId, 'Champs requis.', true);
        }
        try {
            let result;
            const url = currentEditData.hotel ? `/admin/hotels/${currentEditData.hotel.id}` : '/admin/hotels';
            const method = currentEditData.hotel ? 'PUT' : 'POST';
            result = await apiRequest(url, method, formData, true);
            displayFeedback(feedbackId, currentEditData.hotel ? 'Hôtel màj.' : `Hôtel ajouté (ID: ${result?.hotelId ?? 'N/A'}).`, false);
            adminHotelForm.reset();
            adminHotelForm.style.display = 'none';
            loadAdminHotels(hotelFilterChain.value || null);
            loadAdminHotelChains();
        } catch (error) {
            displayFeedback(feedbackId, `Erreur sauvegarde: ${error.message}`, true);
        }
    });
    adminHotelsListDiv?.addEventListener('click', async (e) => {
        if (e.target.classList.contains('edit-button')) {
            const id = e.target.dataset.id;
            try {
                const hotelData = await apiRequest(`/admin/hotels/${id}`, 'GET', null, true);
                currentEditData.hotel = hotelData;
                adminHotelForm.elements['hotel-edit-id'].value = hotelData.id;
                adminHotelForm.elements['chainId'].value = hotelData.chainId;
                adminHotelForm.elements['name'].value = hotelData.name;
                adminHotelForm.elements['address'].value = hotelData.address;
                adminHotelForm.elements['city'].value = hotelData.city;
                adminHotelForm.elements['email'].value = hotelData.email;
                adminHotelForm.elements['telephone'].value = hotelData.telephone;
                adminHotelForm.elements['rating'].value = hotelData.rating;
                adminHotelForm.style.display = 'block';
                window.scrollTo({ top: adminHotelForm.offsetTop - 80, behavior: 'smooth' });
            } catch(error) {
                displayFeedback('admin-hotel-feedback', `Erreur chargement détails hôtel: ${error.message}`, true);
            }
        } else if (e.target.classList.contains('delete-button')) {
            const id = e.target.dataset.id;
            if (confirm(`Supprimer hôtel ID ${id} ? (Supprime aussi chambres)`)) {
                await deleteAdminItem(`/admin/hotels/${id}`, 'admin-hotel-feedback', () => {
                    loadAdminHotels(hotelFilterChain.value || null);
                    loadAdminHotelChains();
                });
            }
        }
    });

    // --- Admin: Room Management ---
    const roomFilterHotel = document.getElementById('admin-room-filter-hotel');
    const roomAmenitiesChecklist = document.getElementById('room-amenities-checklist');
    const roomViewsChecklist = document.getElementById('room-views-checklist');
    let allAmenities = []; 
    let allViewTypes = []; 
    let allRoomTypes = [];

    async function loadRoomFormData() {
        try {
            const [amenities, views, types] = await Promise.all([
                apiRequest('/amenities', 'GET', null, false).catch(e=> {console.error("Amenity load fail",e); return[];}),
                apiRequest('/viewtypes', 'GET', null, false).catch(e=> {console.error("Viewtype load fail",e); return[];}),
                apiRequest('/roomtypes', 'GET', null, false).catch(e=> {console.error("Roomtype load fail",e); return[];})
            ]);
            allAmenities = amenities || [];
            allViewTypes = views || [];
            allRoomTypes = types || [];

            // Populate room type select using room type names (per DTO: RoomType is a string)
            populateSelect('room-type', allRoomTypes, 'name', 'name', 'Choisir Type...');

            // Update checklists to use names rather than IDs
            roomAmenitiesChecklist.innerHTML = allAmenities.map(a => `<label class="checkbox-group"><input type="checkbox" name="amenities" value="${a.name}">${a.name}</label>`).join('');
            roomViewsChecklist.innerHTML = allViewTypes.map(v => `<label class="checkbox-group"><input type="checkbox" name="viewTypes" value="${v.name}">${v.name}</label>`).join('');
        } catch (error) {
            console.error("Failed loading room form data:", error);
            displayFeedback('admin-room-feedback', 'Erreur chargement données chambre.', true);
        }
    }

    async function loadAdminRooms(hotelId = null) {
        const feedbackId = 'admin-room-feedback';
        let endpoint = '/admin/rooms';
        if (hotelId) {
            endpoint += `?hotelId=${hotelId}`;
            console.warn("Assumed room filtering endpoint:", endpoint);
        }
        try {
            const rooms = await apiRequest(endpoint, 'GET', null, true);
            adminRoomsListDiv.innerHTML = createHtmlTable(
                rooms || [],
                ['roomId', 'hotelId', 'capacity', 'number', 'floor', 'price', 'isExtensible'],
                'room'
            );
        } catch (error) {
            displayFeedback(feedbackId, `Erreur chargement chambres: ${error.message}`, true);
            adminRoomsListDiv.innerHTML = '<p>Impossible de charger.</p>';
        }
    }
    roomFilterHotel?.addEventListener('change', (e) => { loadAdminRooms(e.target.value || null); });
    document.getElementById('add-room-button')?.addEventListener('click', () => {
        currentEditData.room = null;
        adminRoomForm.reset();
        adminRoomForm.querySelectorAll('input[type="checkbox"]').forEach(cb => cb.checked = false);
        adminRoomForm.style.display = 'block';
        document.getElementById('room-edit-id').value = '';
        if(document.getElementById('room-hotel-select').options.length <= 1) {
            loadAdminHotels().catch(err => console.error("Failed to populate hotel select", err));
        }
    });
    document.getElementById('cancel-room-button')?.addEventListener('click', () => {
        adminRoomForm.style.display = 'none';
        adminRoomForm.reset();
        clearAllFeedback();
    });
    adminRoomForm?.addEventListener('submit', async (e) => {
        e.preventDefault();
        const feedbackId = 'admin-room-feedback';
        clearAllFeedback();
        const selectedAmenities = Array.from(e.target.elements['amenities']).filter(cb => cb.checked).map(cb => cb.value);
        const selectedViewTypes = Array.from(e.target.elements['viewTypes']).filter(cb => cb.checked).map(cb => cb.value);

        const formData = {
            id: currentEditData.room ? parseInt(currentEditData.room.roomId) : undefined,
            hotelId: parseInt(e.target.elements['hotelId'].value),
            capacity: parseInt(e.target.elements['capacity'].value),
            number: e.target.elements['number'].value.trim(),
            floor: e.target.elements['floor'].value.trim(),
            surfaceArea: parseFloat(e.target.elements['surfaceArea'].value),
            price: parseFloat(e.target.elements['price'].value),
            telephone: e.target.elements['telephone'].value.trim(),
            roomType: e.target.elements['roomType'].value, // now a string per DTO
            isExtensible: e.target.elements['isExtensible'].checked,
            amenities: selectedAmenities,
            viewTypes: selectedViewTypes
        };

        if (isNaN(formData.hotelId) || !formData.number || !formData.floor || isNaN(formData.capacity) || isNaN(formData.price) || !formData.telephone || isNaN(formData.surfaceArea)) {
            return displayFeedback(feedbackId, 'Champs requis invalides.', true);
        }

        try {
            let result;
            const url = currentEditData.room ? `/admin/rooms/${currentEditData.room.roomId}` : '/admin/rooms';
            const method = currentEditData.room ? 'PUT' : 'POST';
            result = await apiRequest(url, method, formData, true);
            displayFeedback(feedbackId, currentEditData.room ? 'Chambre màj.' : `Chambre ajoutée (ID: ${result?.roomId ?? 'N/A'}).`, false);
            adminRoomForm.reset();
            adminRoomForm.style.display = 'none';
            loadAdminRooms(roomFilterHotel.value || null);
            loadAdminHotels();
        } catch (error) {
            displayFeedback(feedbackId, `Erreur sauvegarde chambre: ${error.message}`, true);
        }
    });
    adminRoomsListDiv?.addEventListener('click', async (e) => {
        if (e.target.classList.contains('edit-button')) {
            const id = e.target.dataset.id;
            try {
                const roomData = await apiRequest(`/admin/rooms/${id}`, 'GET', null, true);
                currentEditData.room = roomData;
                adminRoomForm.elements['room-edit-id'].value = roomData.roomId;
                adminRoomForm.elements['hotelId'].value = roomData.hotelId;
                adminRoomForm.elements['number'].value = roomData.number;
                adminRoomForm.elements['floor'].value = roomData.floor;
                adminRoomForm.elements['capacity'].value = roomData.capacity;
                adminRoomForm.elements['roomType'].value = roomData.roomType;
                adminRoomForm.elements['price'].value = roomData.price;
                adminRoomForm.elements['telephone'].value = roomData.telephone;
                adminRoomForm.elements['surfaceArea'].value = roomData.surfaceArea;
                adminRoomForm.elements['isExtensible'].checked = roomData.isExtensible;
                // Set checkboxes based on returned arrays
                const roomAmenities = roomData.amenities || [];
                adminRoomForm.querySelectorAll('input[name="amenities"]').forEach(cb => {
                    cb.checked = roomAmenities.includes(cb.value);
                });
                const roomViews = roomData.viewTypes || [];
                adminRoomForm.querySelectorAll('input[name="viewTypes"]').forEach(cb => {
                    cb.checked = roomViews.includes(cb.value);
                });
                adminRoomForm.style.display = 'block';
                window.scrollTo({ top: adminRoomForm.offsetTop - 80, behavior: 'smooth' });
            } catch (error) {
                displayFeedback('admin-room-feedback', `Erreur chargement chambre: ${error.message}`, true);
            }
        } else if (e.target.classList.contains('delete-button')) {
            const id = e.target.dataset.id;
            if (confirm(`Supprimer chambre ID ${id} ?`)) {
                await deleteAdminItem(`/admin/rooms/${id}`, 'admin-room-feedback', () => {
                    loadAdminRooms(roomFilterHotel.value || null);
                    loadAdminHotels();
                });
            }
        }
    });

    // --- Admin: Client Account Management ---
    async function loadAdminClients() {
        const feedbackId = 'admin-client-feedback';
        try {
            const clients = await apiRequest('/admin/accounts/clients', 'GET', null, true);
            adminClientsListDiv.innerHTML = createHtmlTable(
                clients || [],
                ['accountId', 'firstName', 'lastName', 'email', 'role', 'createdAt', 'updatedAt'],
                'client',
                (acc) => ({ ...acc,
                            createdAt: new Date(acc.createdAt).toLocaleDateString(),
                            updatedAt: new Date(acc.updatedAt).toLocaleDateString()
                          })
            );
        } catch (error) {
            displayFeedback(feedbackId, `Erreur chargement clients: ${error.message}`, true);
            adminClientsListDiv.innerHTML = '<p>Impossible de charger.</p>';
        }
    }
    document.getElementById('add-client-acc-button')?.addEventListener('click', () => {
        currentEditData.client = null;
        adminClientForm.reset();
        adminClientForm.style.display = 'block';
        document.getElementById('client-edit-id').value = '';
        adminClientForm.elements['sin'].disabled = false;
    });
    document.getElementById('cancel-client-acc-button')?.addEventListener('click', () => {
        adminClientForm.style.display = 'none';
        adminClientForm.reset();
        clearAllFeedback();
    });
    adminClientForm?.addEventListener('submit', async (e) => {
        e.preventDefault();
        const feedbackId = 'admin-client-feedback';
        clearAllFeedback();
        const isEditing = !!currentEditData.client;
        const formData = {
            sin: isEditing ? undefined : e.target.elements['sin'].value.trim(),
            firstName: e.target.elements['firstName'].value.trim(),
            lastName: e.target.elements['lastName'].value.trim(),
            address: e.target.elements['address'].value.trim() || undefined,
            phone: e.target.elements['phone'].value.trim() || undefined,
            email: e.target.elements['email'].value.trim(),
        };
        if (!isEditing && !formData.sin) return displayFeedback(feedbackId, 'SIN requis pour création.', true);
        if (!formData.firstName || !formData.lastName || !formData.email) return displayFeedback(feedbackId, 'Prénom, Nom, Email requis.', true);
        if (!isEditing && (formData.sin.length !== 9 || !/^\d{9}$/.test(formData.sin))) return displayFeedback(feedbackId, 'SIN doit être 9 chiffres.', true);

        try {
            let result;
            const accountId = currentEditData.client?.accountId;
            const url = isEditing ? `/admin/accounts/clients/${accountId}` : '/admin/accounts/clients';
            const method = isEditing ? 'PATCH' : 'POST';
            result = await apiRequest(url, method, formData, true);
            displayFeedback(feedbackId, isEditing ? 'Compte Client màj.' : `Compte Client créé (ID: ${result?.accountId ?? 'N/A'}).`, false);
            adminClientForm.reset();
            adminClientForm.style.display = 'none';
            loadAdminClients();
        } catch (error) {
            displayFeedback(feedbackId, `Erreur sauvegarde client: ${error.message}`, true);
        }
    });
    adminClientsListDiv?.addEventListener('click', async (e) => {
        if (e.target.classList.contains('edit-button')) {
            const accountId = e.target.dataset.id;
            try {
                const clientData = await apiRequest(`/admin/accounts/${accountId}`, 'GET', null, true);
                currentEditData.client = clientData;
                currentEditData.client.accountId = accountId;
                adminClientForm.elements['client-edit-id'].value = accountId;
                adminClientForm.elements['sin'].value = clientData.sin;
                adminClientForm.elements['sin'].disabled = true;
                adminClientForm.elements['firstName'].value = clientData.firstName;
                adminClientForm.elements['lastName'].value = clientData.lastName;
                adminClientForm.elements['address'].value = clientData.address || '';
                adminClientForm.elements['phone'].value = clientData.phone || '';
                adminClientForm.elements['email'].value = clientData.email;
                adminClientForm.style.display = 'block';
                window.scrollTo({ top: adminClientForm.offsetTop - 80, behavior: 'smooth' });
            } catch (error) {
                displayFeedback('admin-client-feedback', `Erreur chargement client: ${error.message}`, true);
            }
        } else if (e.target.classList.contains('delete-button')) {
            const accountId = e.target.dataset.id;
            if (confirm(`Supprimer compte client ID ${accountId} ?`)) {
                await deleteAdminItem(`/admin/accounts/clients/${accountId}`, 'admin-client-feedback', loadAdminClients);
            }
        }
    });

    // --- Admin: Employee Account Management ---
    async function loadAdminEmployees() {
        const feedbackId = 'admin-employee-feedback';
        try {
            const employees = await apiRequest('/admin/accounts/employees', 'GET', null, true);
            adminEmployeesListDiv.innerHTML = createHtmlTable(
                employees || [],
                ['accountId', 'firstName', 'lastName', 'email', 'role'],
                'employee',
                (acc) => ({ ...acc })
            );
        } catch (error) {
            displayFeedback(feedbackId, `Erreur chargement employés: ${error.message}`, true);
            adminEmployeesListDiv.innerHTML = '<p>Impossible de charger.</p>';
        }
    }
    document.getElementById('add-employee-acc-button')?.addEventListener('click', () => {
        currentEditData.employee = null;
        adminEmployeeForm.reset();
        adminEmployeeForm.style.display = 'block';
        document.getElementById('employee-edit-id').value = '';
        adminEmployeeForm.elements['sin'].disabled = false;
        if(document.getElementById('employee-acc-hotel').options.length <= 1) {
            loadAdminHotels().catch(err => console.error("Failed hotel load for emp form", err));
        }
    });
    document.getElementById('cancel-employee-acc-button')?.addEventListener('click', () => {
        adminEmployeeForm.style.display = 'none';
        adminEmployeeForm.reset();
        clearAllFeedback();
    });
    adminEmployeeForm?.addEventListener('submit', async (e) => {
        e.preventDefault();
        const feedbackId = 'admin-employee-feedback';
        clearAllFeedback();
        const isEditing = !!currentEditData.employee;
        const formData = {
            sin: isEditing ? undefined : e.target.elements['sin'].value.trim(),
            firstName: e.target.elements['firstName'].value.trim(),
            lastName: e.target.elements['lastName'].value.trim(),
            address: e.target.elements['address'].value.trim(),
            phone: e.target.elements['phone'].value.trim(),
            email: e.target.elements['email'].value.trim(),
            hotelId: parseInt(e.target.elements['hotelId'].value),
            position: e.target.elements['position'].value.trim(),
        };

        if (!isEditing && !formData.sin) return displayFeedback(feedbackId, 'SIN requis pour création.', true);
        if (!formData.firstName || !formData.lastName || !formData.address || !formData.phone || !formData.email || isNaN(formData.hotelId) || !formData.position) {
            return displayFeedback(feedbackId, 'Champs requis.', true);
        }
        try {
            let result;
            const accountId = currentEditData.employee?.accountId;
            const url = isEditing ? `/admin/accounts/employees/${accountId}` : '/admin/accounts/employees';
            const method = isEditing ? 'PATCH' : 'POST';
            result = await apiRequest(url, method, formData, true);
            displayFeedback(feedbackId, isEditing ? 'Compte Employé màj.' : `Compte Employé créé (ID: ${result?.accountId ?? 'N/A'}).`, false);
            adminEmployeeForm.reset();
            adminEmployeeForm.style.display = 'none';
            loadAdminEmployees();
        } catch (error) {
            displayFeedback(feedbackId, `Erreur sauvegarde employé: ${error.message}`, true);
        }
    });
    adminEmployeesListDiv?.addEventListener('click', async (e) => {
        if (e.target.classList.contains('edit-button')) {
            const accountId = e.target.dataset.id;
            try {
                const empData = await apiRequest(`/admin/accounts/${accountId}`, 'GET', null, true);
                currentEditData.employee = empData;
                currentEditData.employee.accountId = accountId;
                adminEmployeeForm.elements['employee-edit-id'].value = accountId;
                adminEmployeeForm.elements['sin'].value = empData.sin;
                adminEmployeeForm.elements['sin'].disabled = true;
                adminEmployeeForm.elements['firstName'].value = empData.firstName;
                adminEmployeeForm.elements['lastName'].value = empData.lastName;
                adminEmployeeForm.elements['address'].value = empData.address;
                adminEmployeeForm.elements['phone'].value = empData.phone;
                adminEmployeeForm.elements['email'].value = empData.email;
                adminEmployeeForm.elements['hotelId'].value = empData.hotelId;
                adminEmployeeForm.elements['position'].value = empData.position;
                adminEmployeeForm.style.display = 'block';
                window.scrollTo({ top: adminEmployeeForm.offsetTop - 80, behavior: 'smooth' });
            } catch (error) {
                displayFeedback('admin-employee-feedback', `Erreur chargement employé: ${error.message}`, true);
            }
        } else if (e.target.classList.contains('delete-button')) {
            const accountId = e.target.dataset.id;
            if (confirm(`Supprimer compte employé ID ${accountId} ?`)) {
                await deleteAdminItem(`/admin/accounts/employees/${accountId}`, 'admin-employee-feedback', loadAdminEmployees);
            }
        }
    });

    // --- deleteAdminItem Helper ---
    async function deleteAdminItem(endpoint, feedbackId, refreshFunction) {
        clearAllFeedback();
        try {
            await apiRequest(endpoint, 'DELETE', null, true);
            displayFeedback(feedbackId, 'Élément supprimé.', false);
            refreshFunction();
        } catch (error) {
            displayFeedback(feedbackId, `Erreur suppression : ${error.message}`, true);
        }
    }

    // --- createHtmlTable Helper ---
    function createHtmlTable(dataArray, columns, itemType, dataProcessor = (item) => item) {
        if (!dataArray || dataArray.length === 0) return '<p>Aucun élément.</p>';
        const headers = columns || Object.keys(dataArray[0] || {});
        if (headers.length === 0) return '<p>Aucune colonne à afficher.</p>';

        let tableHtml = '<table><thead><tr>';
        headers.forEach(header => {
            const title = header.replace(/([A-Z])/g, ' $1').replace(/^./, str => str.toUpperCase());
            tableHtml += `<th>${title}</th>`;
        });
        tableHtml += '<th>Actions</th></tr></thead><tbody>';

        dataArray.forEach(item => {
            const processedItem = dataProcessor(item);
            const idField = itemType === 'client' || itemType === 'employee' ? 'accountId' : (itemType === 'room' ? 'roomId' : 'id');
            const itemId = processedItem[idField];
            tableHtml += `<tr>`;
            headers.forEach((header) => {
                const title = header.replace(/([A-Z])/g, ' $1').replace(/^./, str => str.toUpperCase());
                let cellValue = processedItem[header];
                if (typeof cellValue === 'object' && cellValue !== null) {
                    cellValue = JSON.stringify(cellValue);
                }
                tableHtml += `<td data-label="${title}">${cellValue ?? 'N/A'}</td>`;
            });
            tableHtml += `<td class="actions" data-label="Actions">
                            <button type="button" class="edit-button" data-id="${itemId}" data-type="${itemType}">Modifier</button>
                            <button type="button" class="delete-button" data-id="${itemId}" data-type="${itemType}">Supprimer</button>
                          </td>`;
            tableHtml += `</tr>`;
        });
        tableHtml += '</tbody></table>';
        return tableHtml;
    }

    // --- Required Views Logic ---
    const viewRoomsPerAreaDiv = document.getElementById('view-rooms-per-area');
    const viewHotelCapacitySelect = document.getElementById('view-hotel-capacity-select');
    const viewHotelCapacityResultDiv = document.getElementById('view-hotel-capacity-result');
    const refreshView1Button = document.getElementById('refresh-view1-button');

    async function loadRequiredViewsData() {
        clearAllFeedback();
        if (localStorage.getItem('role') !== 'employee') return;
        console.log("Loading required views...");
        loadView1_RoomsPerArea();
        if (viewHotelCapacitySelect && viewHotelCapacitySelect.options.length <= 1) {
            loadAdminHotels().catch(err => console.error("Hotel load failed for view 2", err));
        }
    }

    async function loadView1_RoomsPerArea() {
        const feedbackId = 'views-feedback';
        viewRoomsPerAreaDiv.innerHTML = '<p>Chargement Vue 1...</p>';
        try {
            const viewData = await apiRequest('/search/zones/rooms', 'GET', null, true);
            if (!viewData || !Array.isArray(viewData) || viewData.length === 0) {
                viewRoomsPerAreaDiv.innerHTML = '<p>Aucune donnée disponible.</p>';
                return;
            }
            let html = '<table><thead><tr><th>Zone (Ville)</th><th>Nb Chambres Disponibles</th></tr></thead><tbody>';
            html += viewData.map(item => `<tr><td data-label="Zone">${item.area || 'N/A'}</td><td data-label="Nb Chambres">${item.available_rooms_count ?? 'N/A'}</td></tr>`).join('');
            html += '</tbody></table>';
            viewRoomsPerAreaDiv.innerHTML = html;
        } catch (error) {
            displayFeedback(feedbackId, `Erreur Vue 1: ${error.message}`, true);
            viewRoomsPerAreaDiv.innerHTML = '<p>Impossible de charger.</p>';
        }
    }

    refreshView1Button?.addEventListener('click', loadView1_RoomsPerArea);

    viewHotelCapacitySelect?.addEventListener('change', async (e) => {
        const hotelId = e.target.value;
        const feedbackId = 'views-feedback';
        viewHotelCapacityResultDiv.innerHTML = '';
        clearAllFeedback();
        if (!hotelId) return;
        viewHotelCapacityResultDiv.innerHTML = '<p>Chargement...</p>';
        try {
            const result = await apiRequest(`/search/hotels/${hotelId}/room-count`, 'GET', null, true);
            const capacity = result?.total_capacity ?? 'N/A';
            viewHotelCapacityResultDiv.innerHTML = `<p><strong>Capacité totale (Hôtel ID ${hotelId}):</strong> ${capacity}</p>`;
        } catch (error) {
            displayFeedback(feedbackId, `Erreur Vue 2: ${error.message}`, true);
            viewHotelCapacityResultDiv.innerHTML = '<p>Impossible de charger.</p>';
        }
    });

    // --- Initial Setup ---
    function initializeApp() {
        console.log("Initializing Sunflower Booking App...");
        const loggedInViaMagic = handleMagicLink();
        if (!loggedInViaMagic) {
            const token = localStorage.getItem('jwt');
            const role = localStorage.getItem('role');
            if (token) {
                if (role === 'client') showView('search-view');
                else if (role === 'employee') { 
                    if(isAdminUser()) showView('admin-dashboard-view'); 
                    else showView('employee-dashboard-view'); 
                } else showView('login-view');
            } else showView('search-view');
        }
        updateNav();
    }

    initializeApp();
});
