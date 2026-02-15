
import type { DefineComponent, SlotsType } from 'vue'
type IslandComponent<T> = DefineComponent<{}, {refresh: () => Promise<void>}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, SlotsType<{ fallback: { error: unknown } }>> & T

type HydrationStrategies = {
  hydrateOnVisible?: IntersectionObserverInit | true
  hydrateOnIdle?: number | true
  hydrateOnInteraction?: keyof HTMLElementEventMap | Array<keyof HTMLElementEventMap> | true
  hydrateOnMediaQuery?: string
  hydrateAfter?: number
  hydrateWhen?: boolean
  hydrateNever?: true
}
type LazyComponent<T> = DefineComponent<HydrationStrategies, {}, {}, {}, {}, {}, {}, { hydrated: () => void }> & T


export const Badge: typeof import("../app/components/ui/badge/index")['Badge']
export const Button: typeof import("../app/components/ui/button/index")['Button']
export const Card: typeof import("../app/components/ui/card/index")['Card']
export const CardAction: typeof import("../app/components/ui/card/index")['CardAction']
export const CardContent: typeof import("../app/components/ui/card/index")['CardContent']
export const CardDescription: typeof import("../app/components/ui/card/index")['CardDescription']
export const CardFooter: typeof import("../app/components/ui/card/index")['CardFooter']
export const CardHeader: typeof import("../app/components/ui/card/index")['CardHeader']
export const CardTitle: typeof import("../app/components/ui/card/index")['CardTitle']
export const Input: typeof import("../app/components/ui/input/index")['Input']
export const Select: typeof import("../app/components/ui/select/index")['Select']
export const SelectContent: typeof import("../app/components/ui/select/index")['SelectContent']
export const SelectGroup: typeof import("../app/components/ui/select/index")['SelectGroup']
export const SelectItem: typeof import("../app/components/ui/select/index")['SelectItem']
export const SelectItemText: typeof import("../app/components/ui/select/index")['SelectItemText']
export const SelectLabel: typeof import("../app/components/ui/select/index")['SelectLabel']
export const SelectScrollDownButton: typeof import("../app/components/ui/select/index")['SelectScrollDownButton']
export const SelectScrollUpButton: typeof import("../app/components/ui/select/index")['SelectScrollUpButton']
export const SelectSeparator: typeof import("../app/components/ui/select/index")['SelectSeparator']
export const SelectTrigger: typeof import("../app/components/ui/select/index")['SelectTrigger']
export const SelectValue: typeof import("../app/components/ui/select/index")['SelectValue']
export const Separator: typeof import("../app/components/ui/separator/index")['Separator']
export const Table: typeof import("../app/components/ui/table/index")['Table']
export const TableBody: typeof import("../app/components/ui/table/index")['TableBody']
export const TableCaption: typeof import("../app/components/ui/table/index")['TableCaption']
export const TableCell: typeof import("../app/components/ui/table/index")['TableCell']
export const TableEmpty: typeof import("../app/components/ui/table/index")['TableEmpty']
export const TableFooter: typeof import("../app/components/ui/table/index")['TableFooter']
export const TableHead: typeof import("../app/components/ui/table/index")['TableHead']
export const TableHeader: typeof import("../app/components/ui/table/index")['TableHeader']
export const TableRow: typeof import("../app/components/ui/table/index")['TableRow']
export const TableUtils: typeof import("../app/components/ui/table/utils")['default']
export const NuxtWelcome: typeof import("../node_modules/nuxt/dist/app/components/welcome.vue")['default']
export const NuxtLayout: typeof import("../node_modules/nuxt/dist/app/components/nuxt-layout")['default']
export const NuxtErrorBoundary: typeof import("../node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']
export const ClientOnly: typeof import("../node_modules/nuxt/dist/app/components/client-only")['default']
export const DevOnly: typeof import("../node_modules/nuxt/dist/app/components/dev-only")['default']
export const ServerPlaceholder: typeof import("../node_modules/nuxt/dist/app/components/server-placeholder")['default']
export const NuxtLink: typeof import("../node_modules/nuxt/dist/app/components/nuxt-link")['default']
export const NuxtLoadingIndicator: typeof import("../node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']
export const NuxtTime: typeof import("../node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']
export const NuxtRouteAnnouncer: typeof import("../node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']
export const NuxtImg: typeof import("../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtImg']
export const NuxtPicture: typeof import("../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtPicture']
export const NuxtPage: typeof import("../node_modules/nuxt/dist/pages/runtime/page")['default']
export const NoScript: typeof import("../node_modules/nuxt/dist/head/runtime/components")['NoScript']
export const Link: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Link']
export const Base: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Base']
export const Title: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Title']
export const Meta: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Meta']
export const Style: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Style']
export const Head: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Head']
export const Html: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Html']
export const Body: typeof import("../node_modules/nuxt/dist/head/runtime/components")['Body']
export const NuxtIsland: typeof import("../node_modules/nuxt/dist/app/components/nuxt-island")['default']
export const LazyBadge: LazyComponent<typeof import("../app/components/ui/badge/index")['Badge']>
export const LazyButton: LazyComponent<typeof import("../app/components/ui/button/index")['Button']>
export const LazyCard: LazyComponent<typeof import("../app/components/ui/card/index")['Card']>
export const LazyCardAction: LazyComponent<typeof import("../app/components/ui/card/index")['CardAction']>
export const LazyCardContent: LazyComponent<typeof import("../app/components/ui/card/index")['CardContent']>
export const LazyCardDescription: LazyComponent<typeof import("../app/components/ui/card/index")['CardDescription']>
export const LazyCardFooter: LazyComponent<typeof import("../app/components/ui/card/index")['CardFooter']>
export const LazyCardHeader: LazyComponent<typeof import("../app/components/ui/card/index")['CardHeader']>
export const LazyCardTitle: LazyComponent<typeof import("../app/components/ui/card/index")['CardTitle']>
export const LazyInput: LazyComponent<typeof import("../app/components/ui/input/index")['Input']>
export const LazySelect: LazyComponent<typeof import("../app/components/ui/select/index")['Select']>
export const LazySelectContent: LazyComponent<typeof import("../app/components/ui/select/index")['SelectContent']>
export const LazySelectGroup: LazyComponent<typeof import("../app/components/ui/select/index")['SelectGroup']>
export const LazySelectItem: LazyComponent<typeof import("../app/components/ui/select/index")['SelectItem']>
export const LazySelectItemText: LazyComponent<typeof import("../app/components/ui/select/index")['SelectItemText']>
export const LazySelectLabel: LazyComponent<typeof import("../app/components/ui/select/index")['SelectLabel']>
export const LazySelectScrollDownButton: LazyComponent<typeof import("../app/components/ui/select/index")['SelectScrollDownButton']>
export const LazySelectScrollUpButton: LazyComponent<typeof import("../app/components/ui/select/index")['SelectScrollUpButton']>
export const LazySelectSeparator: LazyComponent<typeof import("../app/components/ui/select/index")['SelectSeparator']>
export const LazySelectTrigger: LazyComponent<typeof import("../app/components/ui/select/index")['SelectTrigger']>
export const LazySelectValue: LazyComponent<typeof import("../app/components/ui/select/index")['SelectValue']>
export const LazySeparator: LazyComponent<typeof import("../app/components/ui/separator/index")['Separator']>
export const LazyTable: LazyComponent<typeof import("../app/components/ui/table/index")['Table']>
export const LazyTableBody: LazyComponent<typeof import("../app/components/ui/table/index")['TableBody']>
export const LazyTableCaption: LazyComponent<typeof import("../app/components/ui/table/index")['TableCaption']>
export const LazyTableCell: LazyComponent<typeof import("../app/components/ui/table/index")['TableCell']>
export const LazyTableEmpty: LazyComponent<typeof import("../app/components/ui/table/index")['TableEmpty']>
export const LazyTableFooter: LazyComponent<typeof import("../app/components/ui/table/index")['TableFooter']>
export const LazyTableHead: LazyComponent<typeof import("../app/components/ui/table/index")['TableHead']>
export const LazyTableHeader: LazyComponent<typeof import("../app/components/ui/table/index")['TableHeader']>
export const LazyTableRow: LazyComponent<typeof import("../app/components/ui/table/index")['TableRow']>
export const LazyTableUtils: LazyComponent<typeof import("../app/components/ui/table/utils")['default']>
export const LazyNuxtWelcome: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/welcome.vue")['default']>
export const LazyNuxtLayout: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-layout")['default']>
export const LazyNuxtErrorBoundary: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-error-boundary.vue")['default']>
export const LazyClientOnly: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/client-only")['default']>
export const LazyDevOnly: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/dev-only")['default']>
export const LazyServerPlaceholder: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/server-placeholder")['default']>
export const LazyNuxtLink: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-link")['default']>
export const LazyNuxtLoadingIndicator: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-loading-indicator")['default']>
export const LazyNuxtTime: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-time.vue")['default']>
export const LazyNuxtRouteAnnouncer: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-route-announcer")['default']>
export const LazyNuxtImg: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtImg']>
export const LazyNuxtPicture: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-stubs")['NuxtPicture']>
export const LazyNuxtPage: LazyComponent<typeof import("../node_modules/nuxt/dist/pages/runtime/page")['default']>
export const LazyNoScript: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['NoScript']>
export const LazyLink: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Link']>
export const LazyBase: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Base']>
export const LazyTitle: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Title']>
export const LazyMeta: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Meta']>
export const LazyStyle: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Style']>
export const LazyHead: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Head']>
export const LazyHtml: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Html']>
export const LazyBody: LazyComponent<typeof import("../node_modules/nuxt/dist/head/runtime/components")['Body']>
export const LazyNuxtIsland: LazyComponent<typeof import("../node_modules/nuxt/dist/app/components/nuxt-island")['default']>

export const componentNames: string[]
