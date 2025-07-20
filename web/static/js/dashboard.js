// Dashboard JavaScript functionality
let currentPage = 1;
let currentLimit = 10;
let currentQuery = '';
let currentType = 'all';

// Initialize dashboard on page load
document.addEventListener('DOMContentLoaded', function() {
    loadDashboard();
    loadPopularContent();
    
    // Add enter key support for search
    document.getElementById('searchInput').addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            performSearch();
        }
    });
});

// Load dashboard data
async function loadDashboard() {
    try {
        const response = await fetch('/api/dashboard');
        const data = await response.json();
        
        if (data.success) {
            updateStatistics(data.data.statistics);
            updateProviderCount(data.data.providers.length);
        }
    } catch (error) {
        console.error('Dashboard yüklenirken hata:', error);
    }
}

// Load popular content
async function loadPopularContent() {
    try {
        const response = await fetch('/api/content/popular?limit=6');
        const data = await response.json();
        
        if (data.success) {
            displayPopularContent(data.data);
        }
    } catch (error) {
        console.error('Popüler içerik yüklenirken hata:', error);
    }
}

// Perform search
async function performSearch() {
    const query = document.getElementById('searchInput').value.trim();
    const contentType = document.getElementById('contentType').value;
    
    if (!query) {
        alert('Lütfen bir arama sorgusu girin');
        return;
    }
    
    currentQuery = query;
    currentType = contentType;
    currentPage = 1;
    
    await executeSearch();
}

// Execute search with current parameters
async function executeSearch() {
    showLoading(true);
    
    try {
        const params = new URLSearchParams({
            q: currentQuery,
            type: currentType,
            page: currentPage,
            limit: currentLimit
        });
        
        const response = await fetch(`/api/search?${params}`);
        const data = await response.json();
        
        if (data.success) {
            displaySearchResults(data.data);
            displayPagination(data.data);
        } else {
            showError('Arama sırasında bir hata oluştu');
        }
    } catch (error) {
        console.error('Arama hatası:', error);
        showError('Arama sırasında bir hata oluştu');
    } finally {
        showLoading(false);
    }
}

// Display search results
function displaySearchResults(result) {
    const container = document.getElementById('searchResults');
    container.innerHTML = '';
    
    if (result.contents.length === 0) {
        container.innerHTML = `
            <div class="col-12 text-center">
                <div class="alert alert-info">
                    <i class="fas fa-info-circle"></i> Arama kriterlerinize uygun içerik bulunamadı.
                </div>
            </div>
        `;
        return;
    }
    
    result.contents.forEach(content => {
        const card = createContentCard(content);
        container.appendChild(card);
    });
}

// Create content card
function createContentCard(content) {
    const col = document.createElement('div');
    col.className = 'col-md-6 col-lg-4 mb-4';
    
    const typeIcon = content.type === 'video' ? 'fas fa-video' : 'fas fa-file-alt';
    const typeBadgeClass = content.type === 'video' ? 'bg-danger' : 'bg-primary';
    const scoreClass = content.final_score > 50 ? 'bg-success' : content.final_score > 25 ? 'bg-warning' : 'bg-secondary';
    
    col.innerHTML = `
        <div class="card content-card h-100" onclick="showContentDetail(${content.id})">
            <div class="card-body">
                <div class="d-flex justify-content-between align-items-start mb-2">
                    <span class="badge ${typeBadgeClass} type-badge">
                        <i class="${typeIcon}"></i> ${content.type}
                    </span>
                    <span class="badge ${scoreClass} score-badge">
                        ${content.final_score.toFixed(1)} puan
                    </span>
                </div>
                <h6 class="card-title">${content.title}</h6>
                <p class="card-text text-muted small">${content.description.substring(0, 100)}${content.description.length > 100 ? '...' : ''}</p>
                <div class="mt-auto">
                    <small class="text-muted">
                        <i class="fas fa-calendar"></i> ${formatDate(content.published_at)}
                    </small>
                </div>
            </div>
        </div>
    `;
    
    return col;
}

// Display popular content
function displayPopularContent(contents) {
    const container = document.getElementById('popularContent');
    container.innerHTML = '';
    
    contents.forEach(content => {
        const card = createPopularContentCard(content);
        container.appendChild(card);
    });
}

// Create popular content card
function createPopularContentCard(content) {
    const col = document.createElement('div');
    col.className = 'col-md-4 mb-3';
    
    const typeIcon = content.type === 'video' ? 'fas fa-video' : 'fas fa-file-alt';
    const typeBadgeClass = content.type === 'video' ? 'bg-danger' : 'bg-primary';
    
    col.innerHTML = `
        <div class="card content-card h-100" onclick="showContentDetail(${content.id})">
            <div class="card-body">
                <div class="d-flex justify-content-between align-items-start mb-2">
                    <span class="badge ${typeBadgeClass} type-badge">
                        <i class="${typeIcon}"></i> ${content.type}
                    </span>
                    <span class="badge bg-success score-badge">
                        ${content.final_score.toFixed(1)} puan
                    </span>
                </div>
                <h6 class="card-title">${content.title}</h6>
                <p class="card-text text-muted small">${content.description.substring(0, 80)}${content.description.length > 80 ? '...' : ''}</p>
            </div>
        </div>
    `;
    
    return col;
}

// Display pagination
function displayPagination(result) {
    const container = document.getElementById('pagination');
    container.innerHTML = '';
    
    if (result.total_pages <= 1) return;
    
    const pagination = document.createElement('nav');
    pagination.innerHTML = `
        <ul class="pagination justify-content-center">
            <li class="page-item ${result.page === 1 ? 'disabled' : ''}">
                <a class="page-link" href="#" onclick="changePage(${result.page - 1})">Önceki</a>
            </li>
            ${generatePageNumbers(result.page, result.total_pages)}
            <li class="page-item ${result.page === result.total_pages ? 'disabled' : ''}">
                <a class="page-link" href="#" onclick="changePage(${result.page + 1})">Sonraki</a>
            </li>
        </ul>
    `;
    
    container.appendChild(pagination);
}

// Generate page numbers
function generatePageNumbers(currentPage, totalPages) {
    let pages = '';
    const start = Math.max(1, currentPage - 2);
    const end = Math.min(totalPages, currentPage + 2);
    
    for (let i = start; i <= end; i++) {
        pages += `
            <li class="page-item ${i === currentPage ? 'active' : ''}">
                <a class="page-link" href="#" onclick="changePage(${i})">${i}</a>
            </li>
        `;
    }
    
    return pages;
}

// Change page
async function changePage(page) {
    currentPage = page;
    await executeSearch();
    window.scrollTo({ top: 0, behavior: 'smooth' });
}

// Show content detail modal
async function showContentDetail(contentId) {
    try {
        const response = await fetch(`/api/content/${contentId}`);
        const data = await response.json();
        
        if (data.success) {
            const content = data.data.content;
            const breakdown = data.data.score_breakdown;
            
            document.getElementById('modalTitle').textContent = content.title;
            document.getElementById('modalLink').href = content.url;
            
            const modalBody = document.getElementById('modalBody');
            modalBody.innerHTML = `
                <div class="row">
                    <div class="col-md-8">
                        <h6>Açıklama</h6>
                        <p>${content.description}</p>
                        
                        <h6>Detaylar</h6>
                        <ul class="list-unstyled">
                            <li><strong>Tür:</strong> ${content.type}</li>
                            <li><strong>Provider:</strong> ${content.provider}</li>
                            <li><strong>Yayın Tarihi:</strong> ${formatDate(content.published_at)}</li>
                            <li><strong>Dil:</strong> ${content.language}</li>
                        </ul>
                        
                        ${content.type === 'video' ? `
                            <ul class="list-unstyled">
                                <li><strong>Görüntülenme:</strong> ${content.views.toLocaleString()}</li>
                                <li><strong>Beğeni:</strong> ${content.likes.toLocaleString()}</li>
                                <li><strong>Süre:</strong> ${formatDuration(content.duration)}</li>
                            </ul>
                        ` : `
                            <ul class="list-unstyled">
                                <li><strong>Okuma Süresi:</strong> ${content.reading_time} dakika</li>
                                <li><strong>Tepki:</strong> ${content.reactions.toLocaleString()}</li>
                            </ul>
                        `}
                    </div>
                    <div class="col-md-4">
                        <h6>Puan Detayları</h6>
                        <div class="card">
                            <div class="card-body">
                                <p><strong>Temel Puan:</strong> ${breakdown.base_score.toFixed(2)}</p>
                                <p><strong>Tür Çarpanı:</strong> ${breakdown.type_multiplier.toFixed(2)}</p>
                                <p><strong>Güncellik Puanı:</strong> ${breakdown.freshness_score.toFixed(2)}</p>
                                <p><strong>Etkileşim Puanı:</strong> ${breakdown.engagement_score.toFixed(2)}</p>
                                <hr>
                                <p><strong>Final Puan:</strong> <span class="badge bg-success">${breakdown.final_score.toFixed(2)}</span></p>
                            </div>
                        </div>
                    </div>
                </div>
            `;
            
            const modal = new bootstrap.Modal(document.getElementById('contentModal'));
            modal.show();
        }
    } catch (error) {
        console.error('İçerik detayı yüklenirken hata:', error);
        alert('İçerik detayı yüklenirken bir hata oluştu');
    }
}

// Refresh content from providers
async function refreshContent() {
    try {
        showLoading(true);
        
        const response = await fetch('/api/providers/refresh', {
            method: 'POST'
        });
        
        const data = await response.json();
        
        if (data.success) {
            alert('İçerikler başarıyla yenilendi!');
            loadDashboard();
            loadPopularContent();
        } else {
            alert('İçerik yenileme sırasında bir hata oluştu');
        }
    } catch (error) {
        console.error('İçerik yenileme hatası:', error);
        alert('İçerik yenileme sırasında bir hata oluştu');
    } finally {
        showLoading(false);
    }
}

// Update statistics
function updateStatistics(stats) {
    document.getElementById('totalContent').textContent = stats.total_content || 0;
    document.getElementById('videoCount').textContent = stats.video_count || 0;
    document.getElementById('textCount').textContent = stats.text_count || 0;
}

// Update provider count
function updateProviderCount(count) {
    document.getElementById('providerCount').textContent = count;
}

// Show/hide loading
function showLoading(show) {
    const loading = document.getElementById('loading');
    const results = document.getElementById('searchResults');
    
    if (show) {
        loading.style.display = 'block';
        results.style.display = 'none';
    } else {
        loading.style.display = 'none';
        results.style.display = 'flex';
    }
}

// Show error message
function showError(message) {
    const container = document.getElementById('searchResults');
    container.innerHTML = `
        <div class="col-12 text-center">
            <div class="alert alert-danger">
                <i class="fas fa-exclamation-triangle"></i> ${message}
            </div>
        </div>
    `;
}

// Format date
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('tr-TR', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
}

// Format duration (seconds to MM:SS)
function formatDuration(seconds) {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
} 