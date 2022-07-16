package destiny

import (
	"archive/zip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	DisplayProperties struct {
		Description string `json:"description"`
		Name        string `json:"name"`
		Icon        string `json:"icon"`
		HasIcon     bool   `json:"hasIcon"`
	} `json:"displayProperties"`
	CollectibleHash      int64  `json:"collectibleHash"`
	IconWatermark        string `json:"iconWatermark"`
	IconWatermarkShelved string `json:"iconWatermarkShelved"`

	Screenshot                 string `json:"screenshot"`
	ItemTypeDisplayName        string `json:"itemTypeDisplayName"`
	FlavorText                 string `json:"flavorText"`
	UIItemDisplayStyle         string `json:"uiItemDisplayStyle"`
	ItemTypeAndTierDisplayName string `json:"itemTypeAndTierDisplayName"`
	DisplaySource              string `json:"displaySource"`
	Action                     struct {
		VerbName                string        `json:"verbName"`
		VerbDescription         string        `json:"verbDescription"`
		IsPositive              bool          `json:"isPositive"`
		RequiredCooldownSeconds int           `json:"requiredCooldownSeconds"`
		RequiredItems           []interface{} `json:"requiredItems"`
		ProgressionRewards      []interface{} `json:"progressionRewards"`
		ActionTypeLabel         string        `json:"actionTypeLabel"`
		RewardSheetHash         int           `json:"rewardSheetHash"`
		RewardItemHash          int           `json:"rewardItemHash"`
		RewardSiteHash          int           `json:"rewardSiteHash"`
		RequiredCooldownHash    int           `json:"requiredCooldownHash"`
		DeleteOnAction          bool          `json:"deleteOnAction"`
		ConsumeEntireStack      bool          `json:"consumeEntireStack"`
		UseOnAcquire            bool          `json:"useOnAcquire"`
	} `json:"action"`
	Inventory struct {
		MaxStackSize                             int    `json:"maxStackSize"`
		BucketTypeHash                           int64  `json:"bucketTypeHash"`
		RecoveryBucketTypeHash                   int    `json:"recoveryBucketTypeHash"`
		TierTypeHash                             int64  `json:"tierTypeHash"`
		IsInstanceItem                           bool   `json:"isInstanceItem"`
		NonTransferrableOriginal                 bool   `json:"nonTransferrableOriginal"`
		TierTypeName                             string `json:"tierTypeName"`
		TierType                                 int    `json:"tierType"`
		ExpirationTooltip                        string `json:"expirationTooltip"`
		ExpiredInActivityMessage                 string `json:"expiredInActivityMessage"`
		ExpiredInOrbitMessage                    string `json:"expiredInOrbitMessage"`
		SuppressExpirationWhenObjectivesComplete bool   `json:"suppressExpirationWhenObjectivesComplete"`
	} `json:"inventory"`
	Stats struct {
		DisablePrimaryStatDisplay bool        `json:"disablePrimaryStatDisplay"`
		StatGroupHash             int         `json:"statGroupHash"`
		Stats                     interface{} `json:"stats"`
		HasDisplayableStats       bool        `json:"hasDisplayableStats"`
		PrimaryBaseStatHash       int         `json:"primaryBaseStatHash"`
	} `json:"stats"`
	EquippingBlock struct {
		UniqueLabelHash       int      `json:"uniqueLabelHash"`
		EquipmentSlotTypeHash int64    `json:"equipmentSlotTypeHash"`
		Attributes            int      `json:"attributes"`
		EquippingSoundHash    int      `json:"equippingSoundHash"`
		HornSoundHash         int      `json:"hornSoundHash"`
		AmmoType              int      `json:"ammoType"`
		DisplayStrings        []string `json:"displayStrings"`
	} `json:"equippingBlock"`
	TranslationBlock struct {
		WeaponPatternHash int64 `json:"weaponPatternHash"`
		Arrangements      []struct {
			ClassHash          int `json:"classHash"`
			ArtArrangementHash int `json:"artArrangementHash"`
		} `json:"arrangements"`
		HasGeometry bool `json:"hasGeometry"`
	} `json:"translationBlock"`
	Preview struct {
		ScreenStyle         string `json:"screenStyle"`
		PreviewVendorHash   int    `json:"previewVendorHash"`
		PreviewActionString string `json:"previewActionString"`
	} `json:"preview"`
	Quality struct {
		ItemLevels                      []interface{} `json:"itemLevels"`
		QualityLevel                    int           `json:"qualityLevel"`
		InfusionCategoryName            string        `json:"infusionCategoryName"`
		InfusionCategoryHash            int64         `json:"infusionCategoryHash"`
		InfusionCategoryHashes          []int64       `json:"infusionCategoryHashes"`
		ProgressionLevelRequirementHash int64         `json:"progressionLevelRequirementHash"`
		CurrentVersion                  int           `json:"currentVersion"`
		Versions                        []struct {
			PowerCapHash int64 `json:"powerCapHash"`
		} `json:"versions"`
		DisplayVersionWatermarkIcons []string `json:"displayVersionWatermarkIcons"`
	} `json:"quality"`
	AcquireRewardSiteHash int `json:"acquireRewardSiteHash"`
	AcquireUnlockHash     int `json:"acquireUnlockHash"`
	Sockets               struct {
		Detail        string `json:"detail"`
		SocketEntries []struct {
			SocketTypeHash                        int64         `json:"socketTypeHash"`
			SingleInitialItemHash                 int           `json:"singleInitialItemHash"`
			ReusablePlugItems                     []interface{} `json:"reusablePlugItems"`
			PreventInitializationOnVendorPurchase bool          `json:"preventInitializationOnVendorPurchase"`
			PreventInitializationWhenVersioning   bool          `json:"preventInitializationWhenVersioning"`
			HidePerksInItemTooltip                bool          `json:"hidePerksInItemTooltip"`
			PlugSources                           int           `json:"plugSources"`
			ReusablePlugSetHash                   int           `json:"reusablePlugSetHash,omitempty"`
			OverridesUIAppearance                 bool          `json:"overridesUiAppearance"`
			DefaultVisible                        bool          `json:"defaultVisible"`
			RandomizedPlugSetHash                 int           `json:"randomizedPlugSetHash,omitempty"`
		} `json:"socketEntries"`
		IntrinsicSockets []struct {
			PlugItemHash   int64 `json:"plugItemHash"`
			SocketTypeHash int   `json:"socketTypeHash"`
			DefaultVisible bool  `json:"defaultVisible"`
		} `json:"intrinsicSockets"`
		SocketCategories []struct {
			SocketCategoryHash int64 `json:"socketCategoryHash"`
			SocketIndexes      []int `json:"socketIndexes"`
		} `json:"socketCategories"`
	} `json:"sockets"`
	TalentGrid struct {
		TalentGridHash   int    `json:"talentGridHash"`
		ItemDetailString string `json:"itemDetailString"`
		HudDamageType    int    `json:"hudDamageType"`
	} `json:"talentGrid"`
	InvestmentStats []struct {
		StatTypeHash          int  `json:"statTypeHash"`
		Value                 int  `json:"value"`
		IsConditionallyActive bool `json:"isConditionallyActive"`
	} `json:"investmentStats"`
	Perks []struct {
		RequirementDisplayString string `json:"requirementDisplayString"`
		PerkHash                 int    `json:"perkHash"`
		PerkVisibility           int    `json:"perkVisibility"`
	} `json:"perks"`
	LoreHash                          int      `json:"loreHash"`
	SummaryItemHash                   int64    `json:"summaryItemHash"`
	AllowActions                      bool     `json:"allowActions"`
	DoesPostmasterPullHaveSideEffects bool     `json:"doesPostmasterPullHaveSideEffects"`
	NonTransferrable                  bool     `json:"nonTransferrable"`
	ItemCategoryHashes                []int    `json:"itemCategoryHashes"`
	SpecialItemType                   int      `json:"specialItemType"`
	ItemType                          int      `json:"itemType"`
	ItemSubType                       int      `json:"itemSubType"`
	ClassType                         int      `json:"classType"`
	BreakerType                       int      `json:"breakerType"`
	Equippable                        bool     `json:"equippable"`
	DamageTypeHashes                  []int64  `json:"damageTypeHashes"`
	DamageTypes                       []int    `json:"damageTypes"`
	DefaultDamageType                 int      `json:"defaultDamageType"`
	DefaultDamageTypeHash             int64    `json:"defaultDamageTypeHash"`
	IsWrapper                         bool     `json:"isWrapper"`
	TraitIds                          []string `json:"traitIds"`
	TraitHashes                       []int64  `json:"traitHashes"`
	Hash                              int64    `json:"hash"`
	Index                             int      `json:"index"`
	Redacted                          bool     `json:"redacted"`
	Blacklisted                       bool     `json:"blacklisted"`
}

type ManifestURL struct {
	Response struct {
		Version                  string `json:"version"`
		MobileAssetContentPath   string `json:"mobileAssetContentPath"`
		MobileGearAssetDataBases []struct {
			Version int    `json:"version"`
			Path    string `json:"path"`
		} `json:"mobileGearAssetDataBases"`
		MobileWorldContentPaths struct {
			En    string `json:"en"`
			Fr    string `json:"fr"`
			Es    string `json:"es"`
			EsMx  string `json:"es-mx"`
			De    string `json:"de"`
			It    string `json:"it"`
			Ja    string `json:"ja"`
			PtBr  string `json:"pt-br"`
			Ru    string `json:"ru"`
			Pl    string `json:"pl"`
			Ko    string `json:"ko"`
			ZhCht string `json:"zh-cht"`
			ZhChs string `json:"zh-chs"`
		} `json:"mobileWorldContentPaths"`
		JSONWorldContentPaths struct {
			En    string `json:"en"`
			Fr    string `json:"fr"`
			Es    string `json:"es"`
			EsMx  string `json:"es-mx"`
			De    string `json:"de"`
			It    string `json:"it"`
			Ja    string `json:"ja"`
			PtBr  string `json:"pt-br"`
			Ru    string `json:"ru"`
			Pl    string `json:"pl"`
			Ko    string `json:"ko"`
			ZhCht string `json:"zh-cht"`
			ZhChs string `json:"zh-chs"`
		} `json:"jsonWorldContentPaths"`
		JSONWorldComponentContentPaths struct {
			En struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"en"`
			Fr struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"fr"`
			Es struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"es"`
			EsMx struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"es-mx"`
			De struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"de"`
			It struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"it"`
			Ja struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"ja"`
			PtBr struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"pt-br"`
			Ru struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"ru"`
			Pl struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"pl"`
			Ko struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"ko"`
			ZhCht struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"zh-cht"`
			ZhChs struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"zh-chs"`
		} `json:"jsonWorldComponentContentPaths"`
		MobileClanBannerDatabasePath string `json:"mobileClanBannerDatabasePath"`
		MobileGearCDN                struct {
			Geometry    string `json:"Geometry"`
			Texture     string `json:"Texture"`
			PlateRegion string `json:"PlateRegion"`
			Gear        string `json:"Gear"`
			Shader      string `json:"Shader"`
		} `json:"mobileGearCDN"`
		IconImagePyramidInfo []interface{} `json:"iconImagePyramidInfo"`
	} `json:"Response"`
	ErrorCode       int    `json:"ErrorCode"`
	ThrottleSeconds int    `json:"ThrottleSeconds"`
	ErrorStatus     string `json:"ErrorStatus"`
	Message         string `json:"Message"`
	MessageData     struct {
	} `json:"MessageData"`
}

type AutoGenerated struct {
	Response struct {
		Version                  string `json:"version"`
		MobileAssetContentPath   string `json:"mobileAssetContentPath"`
		MobileGearAssetDataBases []struct {
			Version int    `json:"version"`
			Path    string `json:"path"`
		} `json:"mobileGearAssetDataBases"`
		MobileWorldContentPaths struct {
			En    string `json:"en"`
			Fr    string `json:"fr"`
			Es    string `json:"es"`
			EsMx  string `json:"es-mx"`
			De    string `json:"de"`
			It    string `json:"it"`
			Ja    string `json:"ja"`
			PtBr  string `json:"pt-br"`
			Ru    string `json:"ru"`
			Pl    string `json:"pl"`
			Ko    string `json:"ko"`
			ZhCht string `json:"zh-cht"`
			ZhChs string `json:"zh-chs"`
		} `json:"mobileWorldContentPaths"`
		JSONWorldContentPaths struct {
			En    string `json:"en"`
			Fr    string `json:"fr"`
			Es    string `json:"es"`
			EsMx  string `json:"es-mx"`
			De    string `json:"de"`
			It    string `json:"it"`
			Ja    string `json:"ja"`
			PtBr  string `json:"pt-br"`
			Ru    string `json:"ru"`
			Pl    string `json:"pl"`
			Ko    string `json:"ko"`
			ZhCht string `json:"zh-cht"`
			ZhChs string `json:"zh-chs"`
		} `json:"jsonWorldContentPaths"`
		JSONWorldComponentContentPaths struct {
			En struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"en"`
			Fr struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"fr"`
			Es struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"es"`
			EsMx struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"es-mx"`
			De struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"de"`
			It struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"it"`
			Ja struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"ja"`
			PtBr struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"pt-br"`
			Ru struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"ru"`
			Pl struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"pl"`
			Ko struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"ko"`
			ZhCht struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"zh-cht"`
			ZhChs struct {
				DestinyNodeStepSummaryDefinition                string `json:"DestinyNodeStepSummaryDefinition"`
				DestinyArtDyeChannelDefinition                  string `json:"DestinyArtDyeChannelDefinition"`
				DestinyArtDyeReferenceDefinition                string `json:"DestinyArtDyeReferenceDefinition"`
				DestinyPlaceDefinition                          string `json:"DestinyPlaceDefinition"`
				DestinyActivityDefinition                       string `json:"DestinyActivityDefinition"`
				DestinyActivityTypeDefinition                   string `json:"DestinyActivityTypeDefinition"`
				DestinyClassDefinition                          string `json:"DestinyClassDefinition"`
				DestinyGenderDefinition                         string `json:"DestinyGenderDefinition"`
				DestinyInventoryBucketDefinition                string `json:"DestinyInventoryBucketDefinition"`
				DestinyRaceDefinition                           string `json:"DestinyRaceDefinition"`
				DestinyTalentGridDefinition                     string `json:"DestinyTalentGridDefinition"`
				DestinyUnlockDefinition                         string `json:"DestinyUnlockDefinition"`
				DestinyMaterialRequirementSetDefinition         string `json:"DestinyMaterialRequirementSetDefinition"`
				DestinySandboxPerkDefinition                    string `json:"DestinySandboxPerkDefinition"`
				DestinyStatGroupDefinition                      string `json:"DestinyStatGroupDefinition"`
				DestinyProgressionMappingDefinition             string `json:"DestinyProgressionMappingDefinition"`
				DestinyFactionDefinition                        string `json:"DestinyFactionDefinition"`
				DestinyVendorGroupDefinition                    string `json:"DestinyVendorGroupDefinition"`
				DestinyRewardSourceDefinition                   string `json:"DestinyRewardSourceDefinition"`
				DestinyUnlockValueDefinition                    string `json:"DestinyUnlockValueDefinition"`
				DestinyRewardMappingDefinition                  string `json:"DestinyRewardMappingDefinition"`
				DestinyRewardSheetDefinition                    string `json:"DestinyRewardSheetDefinition"`
				DestinyItemCategoryDefinition                   string `json:"DestinyItemCategoryDefinition"`
				DestinyDamageTypeDefinition                     string `json:"DestinyDamageTypeDefinition"`
				DestinyActivityModeDefinition                   string `json:"DestinyActivityModeDefinition"`
				DestinyMedalTierDefinition                      string `json:"DestinyMedalTierDefinition"`
				DestinyAchievementDefinition                    string `json:"DestinyAchievementDefinition"`
				DestinyActivityGraphDefinition                  string `json:"DestinyActivityGraphDefinition"`
				DestinyActivityInteractableDefinition           string `json:"DestinyActivityInteractableDefinition"`
				DestinyBondDefinition                           string `json:"DestinyBondDefinition"`
				DestinyCharacterCustomizationCategoryDefinition string `json:"DestinyCharacterCustomizationCategoryDefinition"`
				DestinyCharacterCustomizationOptionDefinition   string `json:"DestinyCharacterCustomizationOptionDefinition"`
				DestinyCollectibleDefinition                    string `json:"DestinyCollectibleDefinition"`
				DestinyDestinationDefinition                    string `json:"DestinyDestinationDefinition"`
				DestinyEntitlementOfferDefinition               string `json:"DestinyEntitlementOfferDefinition"`
				DestinyEquipmentSlotDefinition                  string `json:"DestinyEquipmentSlotDefinition"`
				DestinyStatDefinition                           string `json:"DestinyStatDefinition"`
				DestinyInventoryItemDefinition                  string `json:"DestinyInventoryItemDefinition"`
				DestinyInventoryItemLiteDefinition              string `json:"DestinyInventoryItemLiteDefinition"`
				DestinyItemTierTypeDefinition                   string `json:"DestinyItemTierTypeDefinition"`
				DestinyLocationDefinition                       string `json:"DestinyLocationDefinition"`
				DestinyLoreDefinition                           string `json:"DestinyLoreDefinition"`
				DestinyMetricDefinition                         string `json:"DestinyMetricDefinition"`
				DestinyObjectiveDefinition                      string `json:"DestinyObjectiveDefinition"`
				DestinyPlatformBucketMappingDefinition          string `json:"DestinyPlatformBucketMappingDefinition"`
				DestinyPlugSetDefinition                        string `json:"DestinyPlugSetDefinition"`
				DestinyPowerCapDefinition                       string `json:"DestinyPowerCapDefinition"`
				DestinyPresentationNodeDefinition               string `json:"DestinyPresentationNodeDefinition"`
				DestinyProgressionDefinition                    string `json:"DestinyProgressionDefinition"`
				DestinyProgressionLevelRequirementDefinition    string `json:"DestinyProgressionLevelRequirementDefinition"`
				DestinyRecordDefinition                         string `json:"DestinyRecordDefinition"`
				DestinyRewardAdjusterPointerDefinition          string `json:"DestinyRewardAdjusterPointerDefinition"`
				DestinyRewardAdjusterProgressionMapDefinition   string `json:"DestinyRewardAdjusterProgressionMapDefinition"`
				DestinyRewardItemListDefinition                 string `json:"DestinyRewardItemListDefinition"`
				DestinySackRewardItemListDefinition             string `json:"DestinySackRewardItemListDefinition"`
				DestinySandboxPatternDefinition                 string `json:"DestinySandboxPatternDefinition"`
				DestinySeasonDefinition                         string `json:"DestinySeasonDefinition"`
				DestinySeasonPassDefinition                     string `json:"DestinySeasonPassDefinition"`
				DestinySocketCategoryDefinition                 string `json:"DestinySocketCategoryDefinition"`
				DestinySocketTypeDefinition                     string `json:"DestinySocketTypeDefinition"`
				DestinyTraitDefinition                          string `json:"DestinyTraitDefinition"`
				DestinyTraitCategoryDefinition                  string `json:"DestinyTraitCategoryDefinition"`
				DestinyUnlockCountMappingDefinition             string `json:"DestinyUnlockCountMappingDefinition"`
				DestinyUnlockEventDefinition                    string `json:"DestinyUnlockEventDefinition"`
				DestinyUnlockExpressionMappingDefinition        string `json:"DestinyUnlockExpressionMappingDefinition"`
				DestinyVendorDefinition                         string `json:"DestinyVendorDefinition"`
				DestinyMilestoneDefinition                      string `json:"DestinyMilestoneDefinition"`
				DestinyActivityModifierDefinition               string `json:"DestinyActivityModifierDefinition"`
				DestinyReportReasonCategoryDefinition           string `json:"DestinyReportReasonCategoryDefinition"`
				DestinyArtifactDefinition                       string `json:"DestinyArtifactDefinition"`
				DestinyBreakerTypeDefinition                    string `json:"DestinyBreakerTypeDefinition"`
				DestinyChecklistDefinition                      string `json:"DestinyChecklistDefinition"`
				DestinyEnergyTypeDefinition                     string `json:"DestinyEnergyTypeDefinition"`
			} `json:"zh-chs"`
		} `json:"jsonWorldComponentContentPaths"`
		MobileClanBannerDatabasePath string `json:"mobileClanBannerDatabasePath"`
		MobileGearCDN                struct {
			Geometry    string `json:"Geometry"`
			Texture     string `json:"Texture"`
			PlateRegion string `json:"PlateRegion"`
			Gear        string `json:"Gear"`
			Shader      string `json:"Shader"`
		} `json:"mobileGearCDN"`
		IconImagePyramidInfo []interface{} `json:"iconImagePyramidInfo"`
	} `json:"Response"`
	ErrorCode       int    `json:"ErrorCode"`
	ThrottleSeconds int    `json:"ThrottleSeconds"`
	ErrorStatus     string `json:"ErrorStatus"`
	Message         string `json:"Message"`
	MessageData     struct {
	} `json:"MessageData"`
}

func GenerateManifest() {
	client := http.Client{}
	request, error := http.NewRequest("GET", "http://www.bungie.net/Platform/Destiny2/Manifest/", nil)
	if error != nil {
		log.Fatal("Unable to create request")
	}
	request.Header.Add("X-API-Key", "a8d4879a0fe04169aa7c7b782265f964")
	response, error := client.Do(request)

	if error != nil {
		log.Fatal("Unable to send request to bungie")
	}

	defer response.Body.Close()
	body, error := ioutil.ReadAll(response.Body)
	if error != nil {
		log.Fatal("Unable to read mobileworld data")
	}

	var data ManifestURL
	json.Unmarshal(body, &data)
	if error != nil {
		log.Fatal("Unable to read mobileworld data")
	}
	url := "http://www.bungie.net" + data.Response.MobileWorldContentPaths.En

	request2, error := http.NewRequest("GET", string(url), nil)
	if error != nil {
		log.Fatal("Unable to create request")
	}
	request2.Header.Add("X-API-Key", "a8d4879a0fe04169aa7c7b782265f964")
	response2, error := client.Do(request2)
	if error != nil {
		log.Fatal("Request to bungie failed")
	}
	defer response2.Body.Close()
	body2, error := ioutil.ReadAll(response2.Body)

	var manifest AutoGenerated
	json.Unmarshal(body2, &manifest)

	os.MkdirAll("./controllers/destiny/manifest", 0755)
	//writting the data to a zip file
	file, error := os.Create("./controllers/destiny/manifest/manifest.zip")
	if error != nil {
		log.Fatal("Unable to generate manifest.zip")
	}

	defer file.Close()
	file.WriteString(string(body2))

	//extracting it to the a manifest.content file for comucating with sqlite.

	resp, err := zip.OpenReader("./controllers/destiny/manifest/manifest.zip")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range resp.File {

		f, err := file.Open()
		if err != nil {
			log.Fatal("Unable to open file")
		}
		defer func() {
			err := f.Close()
			if err != nil {
				log.Fatal("Unable to close file")
			}
		}()

		path := filepath.Join("./controllers/destiny/manifest", file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), file.Mode())
			fs, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())

			if err != nil {
				log.Fatal("Unable to open the file")
			}

			defer func() {
				err := f.Close()
				if err != nil {
					log.Fatal("Unable to close the file")
				}
			}()

			_, err = io.Copy(fs, f)
			if err != nil {
				panic("Unable to copy data")
			}
		}
	}

	//attempting to do this in a database
	dbfile, error := os.Create("./controllers/destiny/manifest/manifest.db")
	if error != nil {
		log.Fatal(error)
	}
	dbfile.Close()

	newDB, error := sql.Open("sqlite3", "./controllers/destiny/manifest/manifest.db")
	if error != nil {
		log.Fatal(error)
		log.Fatal("Unable to open database")
	}

	_, error = newDB.Exec(
		"DROP TABLE IF EXISTS `DestinyInventoryItemDefinition`;" +
			"CREATE TABLE `DestinyInventoryItemDefinition` (`hash` VARCHAR(30) NOT NULL PRIMARY KEY, `json` BLOB NOT NULL);") //+
	//"CREATE TABLE `DestinySandboxPerkDefinition` (`hash` VARCHAR(30) NOT NULL PRIMARY KEY, `json` BLOB NOT NULL);")
	if error != nil {
		log.Fatal(error)
		log.Fatal("Unable to create table")
	}

	//putting it all into a manifest.content file with the hashes rather than id
	db, error := sql.Open("sqlite3", "./controllers/destiny/manifest/world_sql_content_c1d4ac435e5ce5b3046fe2d0e6190ce4.content")
	if error != nil {
		log.Fatal("Unable to load destiny manifest file")

	}

	rows, error := db.Query("SELECT * FROM DestinyInventoryItemDefinition")
	if error != nil {
		log.Fatal("Unable to query the destiny manifest")
	}

	var idd int
	var jsondata string
	for rows.Next() {
		rows.Scan(&idd, &jsondata)
		var tmp Item
		json.Unmarshal([]byte(jsondata), &tmp)
		hash := fmt.Sprintf("%v", tmp.Hash)

		out, _ := json.Marshal(tmp)

		_, error = newDB.Exec("INSERT INTO DestinyInventoryItemDefinition (hash, json) VALUES (?,?)", hash, string(out))
		if error != nil {

			fmt.Println("Unable to inserts destiny items to manifest")
			log.Fatal(error)
		}

	}

	/*
		rows, error = db.Query("SELECT * FROM DestinySandboxPerkDefinition")
		if error != nil {
			log.Fatal("Unable to query the destiny manifest")
		}

		for rows.Next() {
			rows.Scan(&idd, &jsondata)
			var tmp Item
			json.Unmarshal([]byte(jsondata), &tmp)

			hash := fmt.Sprintf("%v", tmp.Hash)
			out, _ := json.Marshal(tmp)

			_, error = newDB.Exec("INSERT INTO DestinySandboxPerkDefinition (hash, json) VALUES (?,?)", hash, string(out))
			if error != nil {
				fmt.Println("Unable to insert perks to manifest")
				fmt.Println(error)
			}

		}
	*/

	newDB.Close()
	//e := os.Remove("./controllers/destiny/manifest/manifest.zip")
	//if e != nil {
	//	log.Fatal("Unable to delete manifest.zip")
	//}
}
