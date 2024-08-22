<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class ActivityCreditWallet
{
    /**
     *
     * @var ?OrchestrationCreditWalletRequest $data
     */
    #[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('\formance\stack\Models\Shared\OrchestrationCreditWalletRequest')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?OrchestrationCreditWalletRequest $data = null;

    /**
     *
     * @var ?string $id
     */
    #[\JMS\Serializer\Annotation\SerializedName('id')]
    #[\JMS\Serializer\Annotation\SkipWhenEmpty]
    public ?string $id = null;

    /**
     * @param  ?OrchestrationCreditWalletRequest  $data
     * @param  ?string  $id
     */
    public function __construct(?OrchestrationCreditWalletRequest $data = null, ?string $id = null)
    {
        $this->data = $data;
        $this->id = $id;
    }
}