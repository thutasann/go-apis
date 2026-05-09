#!/bin/bash

echo "🔍 Comprehensive System Data Usage Analysis on macOS..."
echo "======================================================="

# Function to format size output
format_size() {
    local path="$1"
    local description="$2"
    if [ -d "$path" ] || [ -f "$path" ]; then
        local size=$(du -sh "$path" 2>/dev/null | cut -f1)
        if [ -n "$size" ]; then
            printf "  %-40s %s\n" "$description:" "$size"
        else
            printf "  %-40s %s\n" "$description:" "Not accessible"
        fi
    else
        printf "  %-40s %s\n" "$description:" "Not found"
    fi
}

# Function to show top largest items in a directory
show_top_items() {
    local path="$1"
    local description="$2"
    local count="${3:-5}"
    
    if [ -d "$path" ]; then
        echo "  📋 Top $count largest items in $description:"
        du -sh "$path"/* 2>/dev/null | sort -hr | head -n "$count" | while read size item; do
            printf "    %-10s %s\n" "$size" "$(basename "$item")"
        done
    fi
}

# 1. Overall Disk Space
echo "📦 DISK SPACE OVERVIEW:"
echo "------------------------"
df -h /
echo ""

# 2. System Root Analysis
echo "🏠 SYSTEM ROOT DIRECTORIES:"
echo "----------------------------"
show_top_items "/" "System Root" 10
echo ""

# 3. User Library Deep Dive
echo "📚 USER LIBRARY ANALYSIS (~$USER/Library):"
echo "--------------------------------------------"
format_size ~/Library "Total Library Size"
echo ""

# Application Support
echo "  🔧 Application Support:"
format_size ~/Library/ApplicationSupport "Application Support Total"
show_top_items ~/Library/ApplicationSupport "Application Support" 10
echo ""

# Caches
echo "  💾 Caches:"
format_size ~/Library/Caches "User Caches Total"
show_top_items ~/Library/Caches "User Caches" 10
echo ""

# Containers
echo "  📦 Containers:"
format_size ~/Library/Containers "Containers Total"
show_top_items ~/Library/Containers "Containers" 10
echo ""

# Group Containers
echo "  🏢 Group Containers:"
format_size ~/Library/GroupContainers "Group Containers Total"
show_top_items ~/Library/GroupContainers "Group Containers" 10
echo ""

# 4. System Library Analysis
echo "🏛️  SYSTEM LIBRARY ANALYSIS (/Library):"
echo "-----------------------------------------"
format_size /Library "System Library Total"
echo ""

format_size /Library/Caches "System Caches"
show_top_items /Library/Caches "System Caches" 5
echo ""

format_size /Library/Application\ Support "System Application Support"
show_top_items /Library/Application\ Support "System Application Support" 5
echo ""

# 5. Time Machine & Snapshots
echo "🕒 TIME MACHINE & SNAPSHOTS:"
echo "-----------------------------"
tmutil listlocalsnapshots / 2>/dev/null | while read snapshot; do
    if [ -n "$snapshot" ]; then
        echo "  📸 $snapshot"
    fi
done || echo "  No local snapshots found"

# Check Time Machine size
format_size /Volumes/*/Backups.backupdb 2>/dev/null "Time Machine Backup Size"
echo ""

# 6. iOS/Device Backups
echo "📱 DEVICE BACKUPS:"
echo "------------------"
format_size ~/Library/Application\ Support/MobileSync "MobileSync Total"
format_size ~/Library/Application\ Support/MobileSync/Backup "iOS Backups"
show_top_items ~/Library/Application\ Support/MobileSync/Backup "iOS Backups" 5
echo ""

# 7. Mail Data
echo "📧 MAIL DATA:"
echo "-------------"
format_size ~/Library/Mail "Mail Data Total"
format_size ~/Library/Application\ Support/AddressBook "Address Book"
format_size ~/Library/Calendars "Calendars"
echo ""

# 8. Development Tools Data
echo "⚒️  DEVELOPMENT DATA:"
echo "--------------------"
format_size ~/.npm "NPM Cache"
format_size ~/.yarn "Yarn Cache"
format_size ~/node_modules "User Node Modules"
format_size ~/.docker "Docker Data"
format_size ~/Library/Developer "Xcode Developer Data"
format_size ~/Library/Application\ Support/Code "VS Code Data"
format_size ~/.vscode "VS Code Extensions"
format_size ~/.cursor "Cursor IDE Data"
echo ""

# 9. Browser Data
echo "🌐 BROWSER DATA:"
echo "----------------"
format_size ~/Library/Application\ Support/Google/Chrome "Chrome Data"
format_size ~/Library/Safari "Safari Data"
format_size ~/Library/Application\ Support/Firefox "Firefox Data"
format_size ~/Library/Caches/com.apple.Safari "Safari Caches"
echo ""

# 10. Downloads and Desktop
echo "📥 USER CONTENT:"
echo "----------------"
format_size ~/Downloads "Downloads Folder"
format_size ~/Desktop "Desktop"
format_size ~/Documents "Documents"
format_size ~/Pictures "Pictures"
format_size ~/Movies "Movies"
format_size ~/Music "Music"
echo ""

# 11. System Logs
echo "📝 SYSTEM LOGS:"
echo "---------------"
format_size /var/log "System Logs (/var/log)"
format_size /Library/Logs "Library Logs"
format_size ~/Library/Logs "User Logs"
show_top_items ~/Library/Logs "User Logs" 5
echo ""

# 12. Trash
echo "🗑️  TRASH:"
echo "----------"
format_size ~/.Trash "Trash Size"
echo ""

# 13. Virtual Memory and Swap
echo "💾 VIRTUAL MEMORY:"
echo "------------------"
format_size /private/var/vm "Virtual Memory Files"
ls -lah /private/var/vm/ 2>/dev/null | grep -E "(swapfile|sleepimage)" || echo "  No swap files found"
echo ""

# 14. System Hidden Files
echo "👻 HIDDEN SYSTEM DATA:"
echo "----------------------"
format_size /.fseventsd "FileSystem Events"
format_size /.DocumentRevisions-V100 "Document Revisions"
format_size /.Spotlight-V100 "Spotlight Index"
format_size /.Trashes "System Trashes"
echo ""

# 15. Large Files Search
echo "🔍 SEARCHING FOR LARGE FILES (>1GB):"
echo "-------------------------------------"
echo "  This may take a moment..."
find ~ -type f -size +1G 2>/dev/null | head -10 | while read file; do
    size=$(ls -lah "$file" | awk '{print $5}')
    printf "  %-10s %s\n" "$size" "$file"
done
echo ""

# 16. Summary Statistics
echo "📊 SUMMARY:"
echo "-----------"
total_disk=$(df -h / | tail -1 | awk '{print $2}')
used_disk=$(df -h / | tail -1 | awk '{print $3}')
available_disk=$(df -h / | tail -1 | awk '{print $4}')
usage_percent=$(df -h / | tail -1 | awk '{print $5}')

echo "  💽 Total Disk Space: $total_disk"
echo "  📊 Used Space: $used_disk ($usage_percent)"
echo "  💚 Available Space: $available_disk"
echo ""

echo "✅ Comprehensive system analysis complete!"
echo "💡 Tip: Run 'sudo du -sh /* | sort -hr' for system-wide analysis (requires admin)"
echo "💡 Tip: Use 'ncdu ~' for interactive directory size explorer"