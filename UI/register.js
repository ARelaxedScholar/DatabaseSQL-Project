// register.js
document.addEventListener('DOMContentLoaded', () => {
    const API_BASE = 'https://sunflower-booking-backend-966219880837.us-central1.run.app';
    const form     = document.getElementById('registerForm');
  
    form.addEventListener('submit', async (e) => {
      e.preventDefault();
  
      // Collect form values
      const sin       = form.sin.value.trim();
      const firstName = form.firstName.value.trim();
      const lastName  = form.lastName.value.trim();
      const address   = form.address.value.trim();
      const phone     = form.phone.value.trim();
      const email     = form.email.value.trim();
      // ISO timestamp for joinDate
      const joinDate  = new Date().toISOString();
  
      // Basic client‑side validation
      if (sin.length !== 9) {
        return alert('Le SIN doit contenir exactement 9 caractères.');
      }
  
      const payload = {
        sin,
        firstName,
        lastName,
        address,
        phone,
        email,
        joinDate
      };
  
      try {
        const res = await fetch(`${API_BASE}/clients/register`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        });
  
        if (!res.ok) {
          // Try to parse error message
          const err = await res.json().catch(() => ({}));
          throw new Error(err.message || 'Erreur lors de l’inscription.');
        }
  
        alert('Inscription réussie ! Vous pouvez maintenant vous connecter.');
        form.reset();
        // Optionnel : rediriger vers la page de connexion
        window.location.href = 'booking.html';
      } catch (err) {
        console.error(err);
        alert(err.message);
      }
    });
  });
  