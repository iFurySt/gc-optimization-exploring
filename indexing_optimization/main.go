/**
 * Package indexing_optimization
 * @Author iFurySt <ifuryst@gmail.com>
 * @Date 2024/7/29
 */

package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type (
	Customer struct {
		ID     int
		Name   string
		Age    int
		Gender string
	}
	Department struct {
		ID   int
		Name string
	}
	CustomerServiceStaff struct {
		ID         int
		Name       string
		Department *Department
		Customers  []*Customer
	}
	FakeManager struct {
		Staff       []*CustomerServiceStaff
		StaffMap    map[int]*CustomerServiceStaff
		Customers   []*Customer
		CustomerMap map[int]*Customer
	}
)

var mgr = &FakeManager{}

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":2112", nil)
	}()

	start(context.Background(), 10*time.Second)

	go func() {
		for {
			staffId := rand.Intn(1e6)
			customers := mgr.GetCustomers(staffId)
			fmt.Println(staffId, len(customers))
			time.Sleep(1 * time.Second)
		}
	}()

	select {}
}

func start(ctx context.Context, interval time.Duration) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				refreshData()
				time.Sleep(interval)
			}
		}
	}()
}

func refreshData() {
	reverse := rand.Intn(2) > 1
	genders := []string{"M", "F", "U"}
	departs := []*Department{
		{ID: 1, Name: "Department1"},
		{ID: 2, Name: "Department2"},
		{ID: 3, Name: "Department3"},
		{ID: 4, Name: "Department4"},
		{ID: 5, Name: "Department5"},
	}
	customers := make([]*Customer, 1e6)
	customerMap := make(map[int]*Customer, 1e6)
	for i := 0; i < 1e6; i++ {
		name := fmt.Sprintf("Customer%d", i)
		if reverse {
			name = fmt.Sprintf("Customer%d", 1e6-i)
		}
		customers[i] = &Customer{
			ID:     i,
			Name:   name,
			Age:    rand.Intn(30),
			Gender: genders[rand.Intn(3)],
		}
		customerMap[i] = customers[i]
	}

	staffMap := make(map[int]*CustomerServiceStaff, 1e6)
	staff := make([]*CustomerServiceStaff, 1e6)
	for i := 0; i < 1e6; i++ {
		name := fmt.Sprintf("Staff%d", i)
		if reverse {
			name = fmt.Sprintf("Staff%d", 1e6-i)
		}
		depart := departs[rand.Intn(5)]
		startIdx := rand.Intn(1e6)
		customers := customers[startIdx : startIdx+rand.Intn(1e6-startIdx)]
		staff[i] = &CustomerServiceStaff{
			ID:         i,
			Name:       name,
			Department: depart,
			Customers:  customers,
		}
		staffMap[i] = staff[i]
	}

	mgr.Customers = customers
	mgr.CustomerMap = customerMap
	mgr.Staff = staff
	mgr.StaffMap = staffMap
}

func (m *FakeManager) GetCustomers(staffId int) []string {
	staff, ok := m.StaffMap[staffId]
	if !ok {
		return nil
	}
	customerNames := make([]string, 0, len(staff.Customers))
	for _, c := range staff.Customers {
		customerNames = append(customerNames, c.Name)
	}
	return customerNames
}
