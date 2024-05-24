package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

type MacHash struct {
	Mac  string
	Hash uint32
	Time string
}

// Function to calculate JOAAT hash
func joaatHash(s string) uint32 {
	var hash uint32 = 0
	for _, c := range s {
		hash += uint32(c)
		hash += (hash << 10)
		hash ^= (hash >> 6)
	}
	hash += (hash << 3)
	hash ^= (hash >> 11)
	hash += (hash << 15)
	return hash
}

// Function to generate random MAC address
func randomMAC() string {
	var mac []string
	for i := 0; i < 6; i++ {
		mac = append(mac, fmt.Sprintf("%02x", rand.Intn(256)))
	}
	return strings.Join(mac, ":")
}

// Function to read existing entries from file
func readExistingEntries(filename string) ([]MacHash, error) {
	var entries []MacHash
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return entries, nil // File does not exist, return empty slice
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var entry MacHash
		_, err := fmt.Sscanf(line, "%s - MAC: %s, JOAAT Hash: 0x%x", &entry.Time, &entry.Mac, &entry.Hash)
		if err == nil {
			entries = append(entries, entry)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Read existing entries from file
	existingEntries, err := readExistingEntries("mac_addresses_sorted.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Map to store unique MAC addresses
	uniqueMACs := make(map[string]bool)
	for _, entry := range existingEntries {
		uniqueMACs[entry.Mac] = true
	}

	// Infinite loop to bruteforce indefinitely
	for {
		// Generate random MAC address
		mac := randomMAC()

		// Check if the MAC address already exists
		if _, exists := uniqueMACs[mac]; exists {
			continue // Skip this iteration if the MAC address is a duplicate
		}

		// Calculate JOAAT hash
		hash := joaatHash(mac)

		// Store MAC address in uniqueMACs
		uniqueMACs[mac] = true

		// Store MAC address, hash, and current time in results slice
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		newEntry := MacHash{Mac: mac, Hash: hash, Time: currentTime}

		// Append the new entry to existing entries
		existingEntries = append(existingEntries, newEntry)

		// Sort all results based on JOAAT hashes in ascending order
		sort.Slice(existingEntries, func(i, j int) bool {
			return existingEntries[i].Hash < existingEntries[j].Hash
		})

		// Open the output file in write mode (truncating the content)
		outFile, err := os.OpenFile("mac_addresses_sorted.txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()

		// Write sorted results to file
		for _, entry := range existingEntries {
			line := fmt.Sprintf("%s - MAC: %s, JOAAT Hash: 0x%08x\n", entry.Time, entry.Mac, entry.Hash)
			if _, err := outFile.WriteString(line); err != nil {
				log.Fatal(err)
			}
		}

		// Print confirmation message
		fmt.Println("Results written to mac_addresses_sorted.txt")

		// Sleep for a short duration to avoid overwhelming the system
		time.Sleep(time.Second)

		// Write the lowest values to lowest.txt
		lowestFile, err := os.OpenFile("lowest.txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer lowestFile.Close()

		// Assuming you want the 10 lowest values
		n := 10
		if len(existingEntries) < 10 {
			n = len(existingEntries)
		}

		for i := 0; i < n; i++ {
			entry := existingEntries[i]
			line := fmt.Sprintf("%s - MAC: %s, JOAAT Hash: 0x%08x\n", entry.Time, entry.Mac, entry.Hash)
			if _, err := lowestFile.WriteString(line); err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("Lowest values written to lowest.txt")
	}
}
