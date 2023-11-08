package database_test

import (
	"context"

	"time"

	"github.com/benstev/opener_common/services/database"
	"github.com/benstev/opener_common/test/fakers"
	"github.com/gobuffalo/envy"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Firebase", func() {

	var db *database.DbService
	var mgr *fakers.Manager

	BeforeEach(func() {
		// fakers.LoadEnv()

		Fb_Credentials_File := envy.Get("FB_CREDENTIALS_FILE", "../../../firebase-key.json")
		ProjectID := envy.Get("FB_PROJECT_ID", "controle-volets-ec809")

		fb, err := database.NewDb(context.Background(), Fb_Credentials_File, ProjectID)
		if err != nil {
			db = nil
		} else {
			db = fb
		}
		Expect(err).NotTo(HaveOccurred())
		mgr = fakers.NewManager()
	})

	It("can listen to gates", func(ctx SpecContext) {

		db.ListenToGate(mgr)
		Eventually(mgr.NGates).WithTimeout(5 * time.Second).WithPolling(1 * time.Second).Should(BeNumerically(">", 0))
	})

	It("can add and remove a phone to a gate", func(ctx SpecContext) {
		db.ListenToGate(mgr)

		gateId := "test_gate_gatephone"
		phoneKey := int64(888888)

		err := db.CreateGate(gateId, &database.Gate{Id: gateId, Name: "Gate for tests"})
		Expect(err).NotTo(HaveOccurred())
		Eventually(mgr.FindGateById).WithArguments(gateId).WithTimeout(2 * time.Second).ShouldNot(BeNil())

		err = db.AddGatePhone(gateId, phoneKey)
		Expect(err).NotTo(HaveOccurred())
		Eventually(mgr.GateUpdated).WithTimeout(2 * time.Second).Should(BeTrue())
		g := mgr.FindGateById(gateId)
		Expect(g).NotTo(BeNil())
		Expect(g.Phones).To(ContainElement(phoneKey))

		err = db.RemoveGatePhone(gateId, phoneKey)
		Expect(err).NotTo(HaveOccurred())
		Eventually(mgr.GateUpdated).WithTimeout(2 * time.Second).Should(BeTrue())
		g = mgr.FindGateById(gateId)
		Expect(g).NotTo(BeNil())
		Expect(g.Phones).NotTo(ContainElement(phoneKey))

		err = db.DeleteGate(gateId)
		Expect(err).NotTo(HaveOccurred())
	})

	It("can set the rolling code to a gate", func(ctx SpecContext) {
		db.ListenToGate(mgr)

		gateId := "test_gate_rollingcode"
		rollingCode := uint32(1234)

		err := db.CreateGate(gateId, &database.Gate{Id: gateId, Name: "Gate for tests"})
		Expect(err).NotTo(HaveOccurred())
		Eventually(mgr.FindGateById).WithArguments(gateId).WithTimeout(2 * time.Second).ShouldNot(BeNil())

		err = db.SetRollingCode(gateId, rollingCode)
		Expect(err).NotTo(HaveOccurred())
		Eventually(mgr.GateUpdated).WithTimeout(2 * time.Second).Should(BeTrue())
		g := mgr.FindGateById(gateId)
		Expect(g).NotTo(BeNil())
		Expect(g.RollingCode).To(Equal(rollingCode))

		err = db.DeleteGate(gateId)
		Expect(err).NotTo(HaveOccurred())
	})

	It("can listen to phones", func(ctx SpecContext) {
		db.ListenToPhone(mgr)
		Eventually(mgr.NPhones).WithTimeout(5 * time.Second).WithPolling(1 * time.Second).Should(BeNumerically(">", 0))
	})

	It("can create, update and delete a phone", func(ctx SpecContext) {
		db.ListenToPhone(mgr)

		phoneKey := int64(99990)
		phoneName := "wakajavaka"

		err := db.CreatePhone(&database.PhoneCreatettributes{Key: phoneKey, FriendlyName: "test_phone", Owner: "beber", Number: "999", Email: "email"})
		Expect(err).NotTo(HaveOccurred())
		Eventually(mgr.FindPhoneByKey).WithArguments(phoneKey).WithTimeout(2 * time.Second).ShouldNot(BeNil())

		p := mgr.FindPhoneByKey(phoneKey)
		Expect(p).NotTo(BeNil())

		// db.UpdatePhone(p.Id, map[string]interface{}{"friendlyName": "wakajavaka"})
		db.UpdatePhone(p.Id, &database.PhoneUpdateAttributes{FriendlyName: &phoneName})
		Eventually(mgr.CheckPhoneFiendlyName).WithArguments(phoneKey, "wakajavaka").WithTimeout(2 * time.Second).Should(BeTrue())

		db.DeletePhone(p.Id)
		Eventually(mgr.FindPhoneByKey).WithArguments(phoneKey).WithTimeout(2 * time.Second).Should(BeNil())
	})

	It("can set activation to a phone", func(ctx SpecContext) {
		db.ListenToPhone(mgr)

		phoneKey := int64(99992)

		err := db.CreatePhone(&database.PhoneCreatettributes{Key: phoneKey, FriendlyName: "test_phone", Owner: "beber", Number: "999", Email: "email"})
		Expect(err).NotTo(HaveOccurred())
		Eventually(mgr.FindPhoneByKey).WithArguments(phoneKey).WithTimeout(2 * time.Second).ShouldNot(BeNil())

		p := mgr.FindPhoneByKey(phoneKey)
		Expect(p).NotTo(BeNil())

		db.SetActivation(p.Id, &database.ActivationRec{Code: 999})
		Eventually(mgr.CheckPhoneHasActivationRecord).WithArguments(phoneKey, 999).WithTimeout(2 * time.Second).Should(BeTrue())

		db.DeletePhone(p.Id)
		Eventually(mgr.FindPhoneByKey).WithArguments(phoneKey).WithTimeout(2 * time.Second).Should(BeNil())
	})

})
